package infrastructure

import (
    "context"

    "emperror.dev/errors"
    "github.com/go-playground/validator"
    "github.com/jackc/pgx/v4"
    "github.com/jackc/pgx/v4/log/zapadapter"
    "go.uber.org/zap"

    "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/core/serializer/json"
    gormPostgres "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/gorm_postgres"
    "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/grpc"
    "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
    otelMetrics "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/otel/metrics"
    "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/otel/tracing"
    postgres "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/postgres_pgx"

    "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/write_service/config"
    "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/write_service/internal/shared/contracts"
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

	meter, err := otelMetrics.AddOtelMetrics(ctx, ic.cfg.OTelMetricsConfig, ic.log)
	if err != nil {
		return nil, nil, err
	}
	infrastructure.Metrics = meter

	var cleanup []func()

	grpcClient, err := grpc.NewGrpcClient(ic.cfg.GRPC)
	if err != nil {
		return nil, nil, err
	}
	cleanup = append(cleanup, func() {
		_ = grpcClient.Close()
	})
	infrastructure.GrpcClient = grpcClient

	gorm, err := gormPostgres.NewGorm(ic.cfg.GormPostgres)
	if err != nil {
		return nil, nil, err
	}
	infrastructure.Gorm = gorm

	pgxConn, err := postgres.NewPgxPoolConn(ctx, ic.cfg.Postgresql, zapadapter.NewLogger(zap.L()), pgx.LogLevelInfo)
	if err != nil {
		return nil, nil, errors.WrapIf(err, "postgresql.NewPgxConn")
	}
	ic.log.Infof("postgres connected: %v", pgxConn.ConnPool.Stat().TotalConns())
	cleanup = append(cleanup, func() {
		pgxConn.Close()
	})
	infrastructure.Pgx = pgxConn

	traceProvider, err := tracing.AddOtelTracing(ic.cfg.OTel)
	if err != nil {
		return nil, nil, err
	}
	cleanup = append(cleanup, func() {
		_ = traceProvider.Shutdown(ctx)
	})

	infrastructure.EventSerializer = json.NewJsonEventSerializer()

	return infrastructure, func() {
		for _, c := range cleanup {
			c()
		}
	}, nil
}
