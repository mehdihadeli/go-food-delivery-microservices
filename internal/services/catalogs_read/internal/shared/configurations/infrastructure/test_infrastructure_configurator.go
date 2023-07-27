package infrastructure

import (
	"context"
	"testing"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/go-playground/validator"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/core/serializer/json"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/grpc"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
	otelMetrics "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/otel/metrics"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/otel/tracing"
	mongoTestContainer "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/test/containers/testcontainer/mongo"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/read_service/config"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/read_service/internal/products/mocks/testData"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/read_service/internal/shared/contracts"
)

type testInfrastructureConfigurator struct {
	log logger.Logger
	cfg *config.Config
	t   *testing.T
}

func NewTestInfrastructureConfigurator(t *testing.T, log logger.Logger, cfg *config.Config) contracts.InfrastructureConfigurator {
	return &testInfrastructureConfigurator{log: log, cfg: cfg, t: t}
}

func (ic *testInfrastructureConfigurator) ConfigInfrastructures(ctx context.Context) (*contracts.InfrastructureConfigurations, func(), error) {
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

	mongo, err := mongoTestContainer.NewMongoTestContainers().Start(ctx, ic.t)
	if err != nil {
		return nil, nil, err
	}
	infrastructure.MongoClient = mongo

	err = seedMongoAndMigration(mongo)
	if err != nil {
		return nil, nil, err
	}

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

func seedMongoAndMigration(db *mongo.Client) error {
	// https://github.com/go-testfixtures/testfixtures#templating
	// seed data
	seedProducts := testData.Products

	//// https://go.dev/doc/faq#convert_slice_of_interface
	data := make([]interface{}, len(seedProducts))
	for i, v := range seedProducts {
		data[i] = v
	}

	collection := db.Database("catalogs_write").Collection("products")
	_, err := collection.InsertMany(context.Background(), data, &options.InsertManyOptions{})
	if err != nil {
		return err
	}

	return nil
}
