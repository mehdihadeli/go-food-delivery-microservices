package infrastructure

import (
	"context"

	"emperror.dev/errors"
	"github.com/go-playground/validator"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/core/serializer/json"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/eventstroredb"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/grpc"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/mongodb"
	otelMetrics "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/otel/metrics"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/otel/tracing"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orders/config"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orders/internal/shared/contracts"
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

	cleanup := []func(){}

	traceProvider, err := tracing.AddOtelTracing(ic.cfg.OTel)
	if err != nil {
		return nil, nil, err
	}
	cleanup = append(cleanup, func() {
		_ = traceProvider.Shutdown(ctx)
	})

	grpcClient, err := grpc.NewGrpcClient(ic.cfg.GRPC)
	if err != nil {
		return nil, nil, err
	}
	cleanup = append(cleanup, func() {
		_ = grpcClient.Close()
	})
	infrastructure.GrpcClient = grpcClient

	mongo, err := mongodb.NewMongoDB(ctx, ic.cfg.Mongo)
	if err != nil {
		return nil, nil, errors.WrapIf(err, "NewMongoDBConn")
	}
	ic.log.Infof("(Mongo connected) SessionsInProgress: {%v}", mongo.NumberSessionsInProgress())
	cleanup = append(cleanup, func() {
		_ = mongo.Disconnect(ctx) // nolint: errcheck
	})
	infrastructure.MongoClient = mongo

	esdb, err := eventstroredb.NewEventStoreDB(ic.cfg.EventStoreConfig)
	if err != nil {
		return nil, nil, err
	}
	esdbSerializer := eventstroredb.NewEsdbSerializer(json.NewJsonMetadataSerializer(), json.NewJsonEventSerializer())
	subscriptionRepository := eventstroredb.NewEsdbSubscriptionCheckpointRepository(esdb, ic.log, esdbSerializer)
	cleanup = append(cleanup, func() {
		_ = esdb.Close() // nolint: errcheck
	})
	infrastructure.Esdb = esdb
	infrastructure.CheckpointRepository = subscriptionRepository
	infrastructure.EsdbSerializer = esdbSerializer
	infrastructure.EventSerializer = json.NewJsonEventSerializer()

	if err != nil {
		return nil, nil, err
	}

	return infrastructure, func() {
		for _, c := range cleanup {
			c()
		}
	}, nil
}
