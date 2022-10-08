package infrastructure

import (
	"context"
	"github.com/go-playground/validator"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/core/serializer/json"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/grpc"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	otelMetrics "github.com/mehdihadeli/store-golang-microservice-sample/pkg/otel/metrics"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/otel/tracing"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/shared/contracts"
)

type infrastructureConfigurator struct {
	log logger.Logger
	cfg *config.Config
}

func NewInfrastructureConfigurator(log logger.Logger, cfg *config.Config) contracts.InfrastructureConfigurator {
	return &infrastructureConfigurator{log: log, cfg: cfg}
}

func (ic *infrastructureConfigurator) ConfigInfrastructures(ctx context.Context) (contracts.InfrastructureConfigurations, func(), error) {
	infrastructure := &infrastructureConfigurations{cfg: ic.cfg, log: ic.log, validator: validator.New()}

	cleanup := []func(){}

	grpcClient, err := grpc.NewGrpcClient(ic.cfg.GRPC)
	if err != nil {
		return nil, nil, err
	}
	cleanup = append(cleanup, func() {
		_ = grpcClient.Close()
	})
	infrastructure.grpcClient = grpcClient

	traceProvider, err := tracing.AddOtelTracing(ic.cfg.OTel)
	if err != nil {
		return nil, nil, err
	}
	cleanup = append(cleanup, func() {
		_ = traceProvider.Shutdown(ctx)
	})

	meter, err := otelMetrics.AddOtelMetrics(ctx, ic.cfg.OTelMetricsConfig, ic.log)
	if err != nil {
		return nil, nil, err
	}
	infrastructure.metrics = meter

	mongoClient, err, mongoCleanup := ic.configMongo(ctx)
	if err != nil {
		return nil, nil, err
	}
	cleanup = append(cleanup, mongoCleanup)
	infrastructure.mongoClient = mongoClient

	redis, err, redisCleanup := ic.configRedis(ctx)
	if err != nil {
		return nil, nil, err
	}
	cleanup = append(cleanup, redisCleanup)
	infrastructure.redis = redis

	infrastructure.eventSerializer = json.NewJsonEventSerializer()

	return infrastructure, func() {
		for _, c := range cleanup {
			defer c()
		}
	}, nil
}
