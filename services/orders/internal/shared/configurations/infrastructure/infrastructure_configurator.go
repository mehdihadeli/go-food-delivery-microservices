package infrastructure

import (
	"context"
	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/go-playground/validator"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/core/serializer"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/core/serializer/json"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/es/contracts"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/es/contracts/projection"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/eventstroredb"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/grpc"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	messageBus "github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/bus"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/producer"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/otel/tracing"
	postgres "github.com/mehdihadeli/store-golang-microservice-sample/pkg/postgres_pgx"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/bus"
	rabbitmqConfigurations "github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/configurations"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/shared/web/custom_middlewares"
	"go.mongodb.org/mongo-driver/mongo"
	"go.opentelemetry.io/otel/sdk/trace"
)

type InfrastructureConfiguration struct {
	Log                          logger.Logger
	Cfg                          *config.Config
	Validator                    *validator.Validate
	Pgx                          *postgres.Pgx
	Metrics                      *OrdersServiceMetrics
	Esdb                         *esdb.Client
	EsdbSerializer               *eventstroredb.EsdbSerializer
	CheckpointRepository         contracts.SubscriptionCheckpointRepository
	ElasticClient                *elasticsearch.Client
	MongoClient                  *mongo.Client
	GrpcClient                   grpc.GrpcClient
	TraceProvider                *trace.TracerProvider
	CustomMiddlewares            cutomMiddlewares.CustomMiddlewares
	Projections                  []projection.IProjection
	EventSerializer              serializer.EventSerializer
	Producer                     producer.Producer
	RabbitMQConfigurationBuilder rabbitmqConfigurations.RabbitMQConfigurationBuilder
	RabbitMQBus                  messageBus.Bus
}

type InfrastructureConfigurator interface {
	ConfigureInfrastructure() error
}

type infrastructureConfigurator struct {
	log logger.Logger
	cfg *config.Config
}

func NewInfrastructureConfigurator(log logger.Logger, cfg *config.Config) *infrastructureConfigurator {
	return &infrastructureConfigurator{log: log, cfg: cfg}
}

func (ic *infrastructureConfigurator) ConfigInfrastructures(ctx context.Context) (*InfrastructureConfiguration, error, func()) {
	infrastructure := &InfrastructureConfiguration{Cfg: ic.cfg, Log: ic.log, Validator: validator.New()}

	metrics := ic.configCatalogsMetrics()
	infrastructure.Metrics = metrics

	cleanup := []func(){}

	traceProvider, err := tracing.AddOtelTracing(ic.cfg.OTel)
	if err != nil {
		return nil, err, nil
	}
	cleanup = append(cleanup, func() {
		_ = traceProvider.Shutdown(ctx)
	})
	infrastructure.TraceProvider = traceProvider

	mongoClient, err, mongoCleanup := ic.configMongo(ctx)
	if err != nil {
		return nil, err, nil
	}
	cleanup = append(cleanup, mongoCleanup)
	infrastructure.MongoClient = mongoClient

	grpcClient, err := grpc.NewGrpcClient(ic.cfg.GRPC)
	if err != nil {
		return nil, err, nil
	}
	cleanup = append(cleanup, func() {
		_ = grpcClient.Close()
	})
	infrastructure.GrpcClient = grpcClient

	esdb, checkpointRepository, esdbSerializer, err, eventStoreCleanup := ic.configEventStore()
	if err != nil {
		return nil, err, nil
	}
	cleanup = append(cleanup, eventStoreCleanup)
	infrastructure.Esdb = esdb
	infrastructure.CheckpointRepository = checkpointRepository
	infrastructure.EsdbSerializer = esdbSerializer

	infrastructure.EventSerializer = json.NewJsonEventSerializer()

	infrastructure.RabbitMQConfigurationBuilder = rabbitmqConfigurations.NewRabbitMQConfigurationBuilder()
	mqBus, err := bus.NewRabbitMQBus(ctx, ic.cfg.RabbitMQ, func(builder rabbitmqConfigurations.RabbitMQConfigurationBuilder) {
		builder = infrastructure.RabbitMQConfigurationBuilder
	}, infrastructure.EventSerializer, ic.log)
	if err != nil {
		return nil, err, nil
	}

	infrastructure.RabbitMQBus = mqBus
	infrastructure.Producer = mqBus

	cleanup = append(cleanup, func() {
		_ = mqBus.Stop(ctx)
	})

	if err != nil {
		return nil, err, nil
	}

	return infrastructure, nil, func() {
		for _, c := range cleanup {
			c()
		}
	}
}
