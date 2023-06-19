package infrastructure

import (
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/fxapp"
)

type InfrastructureConfigurator struct {
	*fxapp.Application
}

func NewInfrastructureConfigurator(fxapp *fxapp.Application) *InfrastructureConfigurator {
	return &InfrastructureConfigurator{
		Application: fxapp,
	}
}

func (ic *InfrastructureConfigurator) ConfigInfrastructures() {
	ic.ResolveFunc(func() error {
		return nil
	})
}

//
//func (ic *infrastructureConfigurator) ConfigInfrastructures(
//	ctx context.Context,
//) (*contracts.InfrastructureConfigurations, func(), error) {
//	infrastructure := &contracts.InfrastructureConfigurations{
//		Cfg:       ic.cfg,
//		Log:       ic.log,
//		Validator: validator.New(),
//	}
//
//	meter, err := otelMetrics.NewOtelMetrics(ctx, ic.cfg.OTelMetricsConfig, ic.log)
//	if err != nil {
//		return nil, nil, err
//	}
//	infrastructure.Metrics = meter
//
//	cleanup := []func(){}
//
//	traceProvider, err := tracing.NewOtelTracing(ic.cfg.OTel)
//	if err != nil {
//		return nil, nil, err
//	}
//	cleanup = append(cleanup, func() {
//		_ = traceProvider.Shutdown(ctx)
//	})
//
//	grpcClient, err := grpc.NewGrpcClient(ic.cfg.GRPC)
//	if err != nil {
//		return nil, nil, err
//	}
//	cleanup = append(cleanup, func() {
//		_ = grpcClient.Close()
//	})
//	infrastructure.GrpcClient = grpcClient
//
//	mongo, err := mongodb.NewMongoDB(ctx, ic.cfg.Mongo)
//	if err != nil {
//		return nil, nil, errors.WrapIf(err, "NewMongoDBConn")
//	}
//	ic.log.Infof("(Mongo connected) SessionsInProgress: {%v}", mongo.NumberSessionsInProgress())
//	cleanup = append(cleanup, func() {
//		_ = mongo.Disconnect(ctx) // nolint: errcheck
//	})
//	infrastructure.MongoClient = mongo
//
//	esdb, err := eventstroredb.NewEventStoreDB(ic.cfg.EventStoreConfig)
//	if err != nil {
//		return nil, nil, err
//	}
//	esdbSerializer := eventstroredb.NewEsdbSerializer(
//		json.NewDefaultMetadataSerializer(),
//		json.NewEventSerializer(),
//	)
//	subscriptionRepository := eventstroredb.NewEsdbSubscriptionCheckpointRepository(
//		esdb,
//		ic.log,
//		esdbSerializer,
//	)
//	cleanup = append(cleanup, func() {
//		_ = esdb.Close() // nolint: errcheck
//	})
//	infrastructure.Esdb = esdb
//	infrastructure.CheckpointRepository = subscriptionRepository
//	infrastructure.EsdbSerializer = esdbSerializer
//	infrastructure.EventSerializer = json.NewEventSerializer()
//
//	if err != nil {
//		return nil, nil, err
//	}
//
//	return infrastructure, func() {
//		for _, c := range cleanup {
//			c()
//		}
//	}, nil
//}
