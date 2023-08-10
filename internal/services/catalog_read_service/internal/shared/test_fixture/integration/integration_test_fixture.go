package integration

import (
	"context"
	"fmt"
	"testing"

	"emperror.dev/errors"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/fxapp/contracts"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/messaging/bus"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/mongodb"
	config2 "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/rabbitmq/config"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/utils"
	rabbithole "github.com/michaelklishin/rabbit-hole"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/otel/trace"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogreadservice/config"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogreadservice/internal/products/contracts/data"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogreadservice/internal/products/mocks/testData"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogreadservice/internal/products/models"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogreadservice/internal/shared/app/test"
)

type IntegrationTestSharedFixture struct {
	suite.Suite
	Cfg                    *config.Config
	Log                    logger.Logger
	Bus                    bus.Bus
	ProductRepository      data.ProductRepository
	ProductCacheRepository data.ProductCacheRepository
	Container              contracts.Container
	RabbitmqCleaner        *rabbithole.Client
	rabbitmqOptions        *config2.RabbitmqOptions
	MongoOptions           *mongodb.MongoDbOptions
	BaseAddress            string
	mongoClient            *mongo.Client
	Items                  []*models.Product
	Tracer                 trace.Tracer
}

func NewIntegrationTestSharedFixture(t *testing.T) *IntegrationTestSharedFixture {
	result := test.NewTestApp().Run(t)

	// https://github.com/michaelklishin/rabbit-hole
	rmqc, _ := rabbithole.NewClient(
		fmt.Sprintf(result.RabbitmqOptions.RabbitmqHostOptions.HttpEndPoint()),
		result.RabbitmqOptions.RabbitmqHostOptions.UserName,
		result.RabbitmqOptions.RabbitmqHostOptions.Password)

	shared := &IntegrationTestSharedFixture{
		Log:                    result.Logger,
		Container:              result.Container,
		Cfg:                    result.Cfg,
		RabbitmqCleaner:        rmqc,
		ProductRepository:      result.ProductRepository,
		ProductCacheRepository: result.ProductCacheRepository,
		Bus:                    result.Bus,
		rabbitmqOptions:        result.RabbitmqOptions,
		MongoOptions:           result.MongoDbOptions,
		BaseAddress:            result.EchoHttpOptions.BasePathAddress(),
		mongoClient:            result.MongoClient,
		Tracer:                 result.Tracer,
	}

	return shared
}

func (i *IntegrationTestSharedFixture) CleanupRabbitmqData() error {
	// https://github.com/michaelklishin/rabbit-hole
	// Get all queues
	queues, err := i.RabbitmqCleaner.ListQueuesIn(i.rabbitmqOptions.RabbitmqHostOptions.VirtualHost)
	if err != nil {
		return err
	}

	// clear each queue
	for _, queue := range queues {
		_, err = i.RabbitmqCleaner.PurgeQueue(
			i.rabbitmqOptions.RabbitmqHostOptions.VirtualHost,
			queue.Name,
		)
		i.Require().NoError(err)
	}

	return nil
}

func (i *IntegrationTestSharedFixture) CleanupMongoData() {
	collections := []string{"products"}
	err := cleanupCollections(i.mongoClient, collections, i.MongoOptions.Database)
	i.Require().NoError(err)
}

func cleanupCollections(db *mongo.Client, collections []string, databaseName string) error {
	database := db.Database(databaseName)
	ctx := context.Background()

	// Iterate over the collections and delete all collections
	for _, collection := range collections {
		collection := database.Collection(collection)

		err := collection.Drop(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

// //////////////////////// Shared Hooks //////////////////////////////////
func (i *IntegrationTestSharedFixture) SetupSuite() {
}

func (i *IntegrationTestSharedFixture) SetupTest() {
	i.T().Log("SetupTest")

	// seed data in each test
	res, err := seedData(i.mongoClient, i.MongoOptions.Database)
	i.Require().NoError(err)
	i.Items = res
}

func (i *IntegrationTestSharedFixture) TearDownTest() {
	i.T().Log("TearDownTest")

	// cleanup test containers with their hooks
	err := i.CleanupRabbitmqData()
	if err != nil {
		i.Require().NoError(err)
	}

	i.CleanupMongoData()
}

func (i *IntegrationTestSharedFixture) TearDownSuite() {
}

func seedData(db *mongo.Client, databaseName string) ([]*models.Product, error) {
	ctx := context.Background()
	//// https://go.dev/doc/faq#convert_slice_of_interface
	data := make([]interface{}, len(testData.Products))
	for i, v := range testData.Products {
		data[i] = v
	}

	collection := db.Database(databaseName).Collection("products")
	_, err := collection.InsertMany(context.Background(), data, &options.InsertManyOptions{})
	if err != nil {
		return nil, errors.WrapIf(err, "error in seed database")
	}

	result, err := mongodb.Paginate[*models.Product](
		ctx,
		utils.NewListQuery(10, 1),
		collection,
		nil,
	)
	return result.Items, nil
}
