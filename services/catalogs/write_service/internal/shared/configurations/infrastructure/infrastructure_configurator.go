package infrastructure

import (
	"context"
	"emperror.dev/errors"
	"github.com/go-playground/validator"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/log/zapadapter"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/core/serializer/json"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/gormPostgres"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/grpc"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	otelMetrics "github.com/mehdihadeli/store-golang-microservice-sample/pkg/otel/metrics"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/otel/tracing"
	postgres "github.com/mehdihadeli/store-golang-microservice-sample/pkg/postgres_pgx"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/shared/contracts"
	"go.uber.org/zap"
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

	meter, err := otelMetrics.AddOtelMetrics(ctx, ic.cfg.OTelMetricsConfig, ic.log)
	if err != nil {
		return nil, nil, err
	}
	infrastructure.metrics = meter

	var cleanup []func()

	grpcClient, err := grpc.NewGrpcClient(ic.cfg.GRPC)
	if err != nil {
		return nil, nil, err
	}
	cleanup = append(cleanup, func() {
		_ = grpcClient.Close()
	})
	infrastructure.grpcClient = grpcClient

	gorm, err := gormPostgres.NewGorm(ic.cfg.GormPostgres)
	if err != nil {
		return nil, nil, err
	}
	infrastructure.gorm = gorm

	pgxConn, err := postgres.NewPgxPoolConn(ic.cfg.Postgresql, zapadapter.NewLogger(zap.L()), pgx.LogLevelInfo)
	if err != nil {
		return nil, nil, errors.WrapIf(err, "postgresql.NewPgxConn")
	}
	ic.log.Infof("postgres connected: %v", pgxConn.ConnPool.Stat().TotalConns())
	cleanup = append(cleanup, func() {
		pgxConn.Close()
	})
	infrastructure.pgx = pgxConn

	traceProvider, err := tracing.AddOtelTracing(ic.cfg.OTel)
	if err != nil {
		return nil, nil, err
	}
	cleanup = append(cleanup, func() {
		_ = traceProvider.Shutdown(ctx)
	})

	infrastructure.eventSerializer = json.NewJsonEventSerializer()

	return infrastructure, func() {
		for _, c := range cleanup {
			c()
		}
	}, nil
}
