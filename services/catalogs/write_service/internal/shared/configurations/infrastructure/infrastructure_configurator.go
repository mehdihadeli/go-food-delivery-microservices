package infrastructure

import (
	"context"
	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/go-playground/validator"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/core/serializer"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/core/serializer/json"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/gormPostgres"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/grpc"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/consumer"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/producer"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/otel/tracing"
	postgres "github.com/mehdihadeli/store-golang-microservice-sample/pkg/postgres_pgx"
	rabbitmqProducer "github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/producer"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/producer/options"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/types"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/config"
	cutomMiddlewares "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/shared/web/middlewares"
	"go.opentelemetry.io/otel/sdk/trace"
)

type InfrastructureConfiguration struct {
	Log                logger.Logger
	Cfg                *config.Config
	TraceProvider      *trace.TracerProvider
	Validator          *validator.Validate
	Pgx                *postgres.Pgx
	Gorm               *gormPostgres.Gorm
	Metrics            *CatalogsServiceMetrics
	Esdb               *esdb.Client
	ElasticClient      *elasticsearch.Client
	GrpcClient         grpc.GrpcClient
	CustomMiddlewares  cutomMiddlewares.CustomMiddlewares
	RabbitMQConnection types.IConnection
	EventSerializer    serializer.EventSerializer
	Producer           producer.Producer
	Consumers          []consumer.Consumer
}

type InfrastructureConfigurator interface {
	ConfigInfrastructures(ctx context.Context) (*InfrastructureConfiguration, error, func())
}

type infrastructureConfigurator struct {
	log logger.Logger
	cfg *config.Config
}

func NewInfrastructureConfigurator(log logger.Logger, cfg *config.Config) InfrastructureConfigurator {
	return &infrastructureConfigurator{log: log, cfg: cfg}
}

func (ic *infrastructureConfigurator) ConfigInfrastructures(ctx context.Context) (*InfrastructureConfiguration, error, func()) {
	infrastructure := &InfrastructureConfiguration{Cfg: ic.cfg, Log: ic.log, Validator: validator.New()}

	metrics := ic.configCatalogsMetrics()
	infrastructure.Metrics = metrics

	var cleanup []func()

	gorm, err := ic.configGorm()
	if err != nil {
		return nil, err, nil
	}
	infrastructure.Gorm = gorm

	grpcClient, err := grpc.NewGrpcClient(ic.cfg.GRPC)
	if err != nil {
		return nil, err, nil
	}
	cleanup = append(cleanup, func() {
		_ = grpcClient.Close()
	})
	infrastructure.GrpcClient = grpcClient

	traceProvider, err := tracing.AddOtelTracing(ic.cfg.OTel)
	if err != nil {
		return nil, err, nil
	}
	cleanup = append(cleanup, func() {
		_ = traceProvider.Shutdown(ctx)
	})
	infrastructure.TraceProvider = traceProvider

	pgx, err, postgresCleanup := ic.configPostgres()
	if err != nil {
		return nil, err, nil
	}
	cleanup = append(cleanup, postgresCleanup)
	infrastructure.Pgx = pgx

	infrastructure.EventSerializer = json.NewJsonEventSerializer()

	connection, err := types.NewRabbitMQConnection(ctx, ic.cfg.RabbitMQ)
	if err != nil {
		return nil, err, nil
	}
	infrastructure.RabbitMQConnection = connection
	cleanup = append(cleanup, func() {
		_ = connection.Close()
	})

	mqProducer, err := rabbitmqProducer.NewRabbitMQProducer(connection, func(builder *options.RabbitMQProducerOptionsBuilder) {}, ic.log, infrastructure.EventSerializer)
	if err != nil {
		return nil, err, nil
	}
	infrastructure.Producer = mqProducer

	if err != nil {
		return nil, err, nil
	}

	return infrastructure, nil, func() {
		for _, c := range cleanup {
			c()
		}
	}
}
