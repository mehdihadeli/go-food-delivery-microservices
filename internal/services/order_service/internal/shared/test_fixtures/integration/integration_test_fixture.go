package integration

import (
	"context"
	"fmt"
	"testing"

	"emperror.dev/errors"
	"github.com/EventStore/EventStore-Client-Go/esdb"
	rabbithole "github.com/michaelklishin/rabbit-hole"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/es/contracts/store"
	config3 "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/eventstroredb/config"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/fxapp/contracts"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/messaging/bus"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/mongodb"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/rabbitmq/config"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/utils"
	config2 "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orderservice/config"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orderservice/internal/orders/contracts/repositories"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orderservice/internal/orders/mocks/testData"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orderservice/internal/orders/models/orders/aggregate"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orderservice/internal/orders/models/orders/read_models"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orderservice/internal/shared/app/test"
	contracts2 "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orderservice/internal/shared/contracts"
	ordersService "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orderservice/internal/shared/grpc/genproto"
)

const (
	orderCollection = "orders"
)

type IntegrationTestSharedFixture struct {
	suite.Suite
	OrderAggregateStore  store.AggregateStore[*aggregate.Order]
	OrderMongoRepository repositories.OrderMongoRepository
	OrdersMetrics        contracts2.OrdersMetrics
	Cfg                  *config2.Config
	Log                  logger.Logger
	Bus                  bus.Bus
	Container            contracts.Container
	RabbitmqCleaner      *rabbithole.Client
	rabbitmqOptions      *config.RabbitmqOptions
	BaseAddress          string
	mongoClient          *mongo.Client
	esdbClient           *esdb.Client
	MongoDbOptions       *mongodb.MongoDbOptions
	EventStoreDbOptions  *config3.EventStoreDbOptions
	Items                []*read_models.OrderReadModel
	OrdersServiceClient  ordersService.OrdersServiceClient
}

func NewIntegrationTestSharedFixture(t *testing.T) *IntegrationTestSharedFixture {
	result := test.NewTestApp().Run(t)

	// https://github.com/michaelklishin/rabbit-hole
	rmqc, _ := rabbithole.NewClient(
		fmt.Sprintf(result.RabbitmqOptions.RabbitmqHostOptions.HttpEndPoint()),
		result.RabbitmqOptions.RabbitmqHostOptions.UserName,
		result.RabbitmqOptions.RabbitmqHostOptions.Password)

	shared := &IntegrationTestSharedFixture{
		Log:                  result.Logger,
		Container:            result.Container,
		Cfg:                  result.Cfg,
		RabbitmqCleaner:      rmqc,
		OrderMongoRepository: result.OrderMongoRepository,
		OrderAggregateStore:  result.OrderAggregateStore,
		MongoDbOptions:       result.MongoDbOptions,
		EventStoreDbOptions:  result.EventStoreDbOptions,
		mongoClient:          result.MongoClient,
		Bus:                  result.Bus,
		rabbitmqOptions:      result.RabbitmqOptions,
		BaseAddress:          result.EchoHttpOptions.BasePathAddress(),
		OrdersServiceClient:  result.OrdersServiceClient,
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
	collections := []string{orderCollection}
	err := cleanupCollections(i.mongoClient, collections, i.MongoDbOptions.Database)
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
func (i *IntegrationTestSharedFixture) SetupTest() {
	i.T().Log("SetupTest")

	// seed data in each test
	res, err := seedReadModelData(i.mongoClient, i.MongoDbOptions.Database)
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

func seedReadModelData(
	db *mongo.Client,
	databaseName string,
) ([]*read_models.OrderReadModel, error) {
	ctx := context.Background()
	//// https://go.dev/doc/faq#convert_slice_of_interface
	data := make([]interface{}, len(testData.Orders))
	for i, v := range testData.Orders {
		data[i] = v
	}

	collection := db.Database(databaseName).Collection("orders")
	_, err := collection.InsertMany(context.Background(), data, &options.InsertManyOptions{})
	if err != nil {
		return nil, errors.WrapIf(err, "error in seed database")
	}

	result, err := mongodb.Paginate[*read_models.OrderReadModel](
		ctx,
		utils.NewListQuery(10, 1),
		collection,
		nil,
	)
	return result.Items, nil
}
