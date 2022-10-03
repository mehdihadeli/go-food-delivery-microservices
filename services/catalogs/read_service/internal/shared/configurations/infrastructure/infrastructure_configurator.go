package infrastructure

import (
	"context"
	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/go-playground/validator"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/core/serializer"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/core/serializer/json"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/grpc"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/consumer"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/producer"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/otel/tracing"
	rabbitmqProducer "github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/producer"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/producer/options"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/types"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/shared/web/middlewares"
	v7 "github.com/olivere/elastic/v7"
	"go.mongodb.org/mongo-driver/mongo"
	"go.opentelemetry.io/otel/sdk/trace"
	"gorm.io/gorm"
)

type InfrastructureConfigurations struct {
	Log                logger.Logger
	Cfg                *config.Config
	Validator          *validator.Validate
	RabbitMQConnection types.IConnection
	Producer           producer.Producer
	Consumers          []consumer.Consumer
	PgConn             *pgxpool.Pool
	Gorm               *gorm.DB
	Metrics            *CatalogsServiceMetrics
	TraceProvider      *trace.TracerProvider
	Esdb               *esdb.Client
	MongoClient        *mongo.Client
	GrpcClient         grpc.GrpcClient
	ElasticClient      *v7.Client
	Redis              redis.UniversalClient
	MiddlewareManager  cutomMiddlewares.CustomMiddlewares
	EventSerializer    serializer.EventSerializer
}

type InfrastructureConfigurator interface {
	ConfigInfrastructures(ctx context.Context) (*InfrastructureConfigurations, error, func())
}

type infrastructureConfigurator struct {
	log logger.Logger
	cfg *config.Config
}

func NewInfrastructureConfigurator(log logger.Logger, cfg *config.Config) InfrastructureConfigurator {
	return &infrastructureConfigurator{log: log, cfg: cfg}
}

func (ic *infrastructureConfigurator) ConfigInfrastructures(ctx context.Context) (*InfrastructureConfigurations, error, func()) {
	infrastructure := &InfrastructureConfigurations{Cfg: ic.cfg, Log: ic.log, Validator: validator.New(), Consumers: make([]consumer.Consumer, 0)}

	metrics := ic.configCatalogsMetrics()
	infrastructure.Metrics = metrics

	cleanup := []func(){}

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

	mongoClient, err, mongoCleanup := ic.configMongo(ctx)
	if err != nil {
		return nil, err, nil
	}
	cleanup = append(cleanup, mongoCleanup)
	infrastructure.MongoClient = mongoClient

	redis, err, redisCleanup := ic.configRedis(ctx)
	if err != nil {
		return nil, err, nil
	}
	cleanup = append(cleanup, redisCleanup)
	infrastructure.Redis = redis

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
			defer c()
		}
	}
}
