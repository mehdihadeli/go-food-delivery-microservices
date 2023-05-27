package infrastructure

import (
	"context"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/redis"

	"emperror.dev/errors"
	"github.com/go-playground/validator"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/core/serializer/json"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/grpc"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/mongodb"
	otelMetrics "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/otel/metrics"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/otel/tracing"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/read_service/config"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/read_service/internal/shared/contracts"
)

type infrastructureConfigurator struct {
	log logger.Logger
	cfg *config.Config
}

func NewInfrastructureConfigurator(log logger.Logger, cfg *config.Config) contracts.InfrastructureConfigurator {
	return &infrastructureConfigurator{log: log, cfg: cfg}
}

func (ic *infrastructureConfigurator) ConfigInfrastructures(ctx context.Context) (*contracts.InfrastructureConfigurations, func(), error) {
	infrastructure := &contracts.InfrastructureConfigurations{Cfg: ic.cfg, Log: ic.log, Validator: validator.New()}

	cleanup := []func(){}

	grpcClient, err := grpc.NewGrpcClient(ic.cfg.GRPC)
	if err != nil {
		return nil, nil, err
	}
	cleanup = append(cleanup, func() {
		_ = grpcClient.Close()
	})
	infrastructure.GrpcClient = grpcClient

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
	infrastructure.Metrics = meter

	mongo, err := mongodb.NewMongoDB(ctx, ic.cfg.Mongo)
	if err != nil {
		return nil, nil, errors.WrapIf(err, "NewMongoDBConn")
	}
	cleanup = append(cleanup, func() {
		_ = mongo.Disconnect(context.Background()) // nolint: errcheck
	})
	infrastructure.MongoClient = mongo

	redis := redis.NewUniversalRedisClient(ic.cfg.Redis)
	ic.log.Infof("Redis connected: %+v", redis.PoolStats())
	cleanup = append(cleanup, func() {
		_ = redis.Close() // nolint: errcheck
	})
	infrastructure.Redis = redis

	infrastructure.EventSerializer = json.NewJsonEventSerializer()

	return infrastructure, func() {
		for _, c := range cleanup {
			defer c()
		}
	}, nil
}
