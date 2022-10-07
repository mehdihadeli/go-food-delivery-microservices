package infrastructure

import (
	"context"
	"github.com/go-playground/validator"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/core/serializer/json"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/grpc"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/otel/tracing"
	rabbitmqBus "github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/bus"
	rabbitmqConfigurations "github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/configurations"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/configurations/rabbitmq"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/shared/contracts"
)

type infrastructureConfigurator struct {
	log logger.Logger
	cfg *config.Config
}

func NewInfrastructureConfigurator(log logger.Logger, cfg *config.Config) contracts.InfrastructureConfigurator {
	return &infrastructureConfigurator{log: log, cfg: cfg}
}

func (ic *infrastructureConfigurator) ConfigInfrastructures(ctx context.Context) (contracts.InfrastructureConfiguration, error, func()) {
	infrastructure := &infrastructureConfigurations{Cfg: ic.cfg, Log: ic.log, Validator: validator.New()}

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
	infrastructure.RabbitMQBus, _ = rabbitmqBus.NewRabbitMQBus(
		ctx,
		infrastructure.GetCfg().RabbitMQ,
		func(builder rabbitmqConfigurations.RabbitMQConfigurationBuilder) {
			rabbitmq.ConfigRabbitMQ(builder, infrastructure)
		},
		infrastructure.GetEventSerializer(),
		infrastructure.GetLog())

	infrastructure.RabbitMQConfigurationBuilder = rabbitmqConfigurations.NewRabbitMQConfigurationBuilder()

	return infrastructure, nil, func() {
		for _, c := range cleanup {
			defer c()
		}
	}
}
