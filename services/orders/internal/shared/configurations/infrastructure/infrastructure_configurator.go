package infrastructure

import (
	"context"
	"emperror.dev/errors"
	"github.com/go-playground/validator"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/core/serializer/json"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/eventstroredb"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/grpc"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/mongodb"
	otelMetrics "github.com/mehdihadeli/store-golang-microservice-sample/pkg/otel/metrics"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/otel/tracing"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/shared/contracts"
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
	infrastructure.grpcClient = grpcClient

	mongo, err := mongodb.NewMongoDB(ctx, ic.cfg.Mongo)
	if err != nil {
		return nil, nil, errors.WrapIf(err, "NewMongoDBConn")
	}
	ic.log.Infof("(Mongo connected) SessionsInProgress: {%v}", mongo.MongoClient.NumberSessionsInProgress())
	cleanup = append(cleanup, func() {
		_ = mongo.Close() // nolint: errcheck
	})
	infrastructure.mongoClient = mongo.MongoClient

	esdb, err := eventstroredb.NewEventStoreDB(ic.cfg.EventStoreConfig)
	if err != nil {
		return nil, nil, err
	}
	esdbSerializer := eventstroredb.NewEsdbSerializer(json.NewJsonMetadataSerializer(), json.NewJsonEventSerializer())
	subscriptionRepository := eventstroredb.NewEsdbSubscriptionCheckpointRepository(esdb, ic.log, esdbSerializer)
	cleanup = append(cleanup, func() {
		_ = esdb.Close() // nolint: errcheck
	})
	infrastructure.esdb = esdb
	infrastructure.checkpointRepository = subscriptionRepository
	infrastructure.esdbSerializer = esdbSerializer
	infrastructure.eventSerializer = json.NewJsonEventSerializer()

	if err != nil {
		return nil, nil, err
	}

	return infrastructure, func() {
		for _, c := range cleanup {
			c()
		}
	}, nil
}
