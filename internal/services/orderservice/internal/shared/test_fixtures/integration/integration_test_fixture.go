package integration

import (
	"context"
	"testing"

	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/core/messaging/bus"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/es/contracts/store"
	config3 "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/eventstroredb/config"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/fxapp/contracts"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/logger"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/mongodb"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/rabbitmq/config"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/utils"
	config2 "github.com/mehdihadeli/go-food-delivery-microservices/internal/services/orderservice/config"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/services/orderservice/internal/orders/contracts/repositories"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/services/orderservice/internal/orders/models/orders/aggregate"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/services/orderservice/internal/orders/models/orders/read_models"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/services/orderservice/internal/shared/app/test"
	contracts2 "github.com/mehdihadeli/go-food-delivery-microservices/internal/services/orderservice/internal/shared/contracts"
	ordersService "github.com/mehdihadeli/go-food-delivery-microservices/internal/services/orderservice/internal/shared/grpc/genproto"

	"emperror.dev/errors"
	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/brianvoe/gofakeit/v6"
	rabbithole "github.com/michaelklishin/rabbit-hole"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	orderCollection = "orders"
)

type IntegrationTestSharedFixture struct {
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

func NewIntegrationTestSharedFixture(
	t *testing.T,
) *IntegrationTestSharedFixture {
	result := test.NewTestApp().Run(t)

	// https://github.com/michaelklishin/rabbit-hole
	rmqc, err := rabbithole.NewClient(
		result.RabbitmqOptions.RabbitmqHostOptions.HttpEndPoint(),
		result.RabbitmqOptions.RabbitmqHostOptions.UserName,
		result.RabbitmqOptions.RabbitmqHostOptions.Password)
	if err != nil {
		result.Logger.Error(
			errors.WrapIf(err, "error in creating rabbithole client"),
		)
	}
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

func (i *IntegrationTestSharedFixture) SetupTest() {
	i.Log.Info("SetupTest started")

	// seed data in each test
	res, err := seedReadModelData(i.mongoClient, i.MongoDbOptions.Database)
	if err != nil {
		i.Log.Error(errors.WrapIf(err, "error in seeding mongodb data"))
	}
	i.Items = res
}

func (i *IntegrationTestSharedFixture) TearDownTest() {
	i.Log.Info("TearDownTest started")

	// cleanup test containers with their hooks
	if err := i.cleanupRabbitmqData(); err != nil {
		i.Log.Error(errors.WrapIf(err, "error in cleanup rabbitmq data"))
	}

	if err := i.cleanupMongoData(); err != nil {
		i.Log.Error(errors.WrapIf(err, "error in cleanup mongodb data"))
	}
}

func (i *IntegrationTestSharedFixture) cleanupRabbitmqData() error {
	// https://github.com/michaelklishin/rabbit-hole
	// Get all queues
	queues, err := i.RabbitmqCleaner.ListQueuesIn(
		i.rabbitmqOptions.RabbitmqHostOptions.VirtualHost,
	)
	if err != nil {
		return err
	}

	// clear each queue
	for _, queue := range queues {
		_, err = i.RabbitmqCleaner.PurgeQueue(
			i.rabbitmqOptions.RabbitmqHostOptions.VirtualHost,
			queue.Name,
		)
		return err
	}

	return nil
}

func (i *IntegrationTestSharedFixture) cleanupMongoData() error {
	collections := []string{orderCollection}
	err := cleanupCollections(
		i.mongoClient,
		collections,
		i.MongoDbOptions.Database,
	)
	return err
}

func cleanupCollections(
	db *mongo.Client,
	collections []string,
	databaseName string,
) error {
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

func seedReadModelData(
	db *mongo.Client,
	databaseName string,
) ([]*read_models.OrderReadModel, error) {
	ctx := context.Background()

	orders := []*read_models.OrderReadModel{
		{
			Id:              gofakeit.UUID(),
			OrderId:         gofakeit.UUID(),
			ShopItems:       generateShopItems(),
			AccountEmail:    gofakeit.Email(),
			DeliveryAddress: gofakeit.Address().Address,
			CancelReason:    gofakeit.Sentence(5),
			TotalPrice:      gofakeit.Float64Range(10, 100),
			DeliveredTime:   gofakeit.Date(),
			Paid:            gofakeit.Bool(),
			Submitted:       gofakeit.Bool(),
			Completed:       gofakeit.Bool(),
			Canceled:        gofakeit.Bool(),
			PaymentId:       gofakeit.UUID(),
			CreatedAt:       gofakeit.Date(),
			UpdatedAt:       gofakeit.Date(),
		},
		{
			Id:              gofakeit.UUID(),
			OrderId:         gofakeit.UUID(),
			ShopItems:       generateShopItems(),
			AccountEmail:    gofakeit.Email(),
			DeliveryAddress: gofakeit.Address().Address,
			CancelReason:    gofakeit.Sentence(5),
			TotalPrice:      gofakeit.Float64Range(10, 100),
			DeliveredTime:   gofakeit.Date(),
			Paid:            gofakeit.Bool(),
			Submitted:       gofakeit.Bool(),
			Completed:       gofakeit.Bool(),
			Canceled:        gofakeit.Bool(),
			PaymentId:       gofakeit.UUID(),
			CreatedAt:       gofakeit.Date(),
			UpdatedAt:       gofakeit.Date(),
		},
	}

	//// https://go.dev/doc/faq#convert_slice_of_interface
	data := make([]interface{}, len(orders))
	for i, v := range orders {
		data[i] = v
	}

	collection := db.Database(databaseName).Collection("orders")
	_, err := collection.InsertMany(
		context.Background(),
		data,
		&options.InsertManyOptions{},
	)
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

func generateShopItems() []*read_models.ShopItemReadModel {
	var shopItems []*read_models.ShopItemReadModel

	for i := 0; i < 3; i++ {
		shopItem := &read_models.ShopItemReadModel{
			Title:       gofakeit.Word(),
			Description: gofakeit.Sentence(3),
			Quantity:    uint64(gofakeit.UintRange(1, 100)),
			Price:       gofakeit.Float64Range(1, 50),
		}

		shopItems = append(shopItems, shopItem)
	}

	return shopItems
}
