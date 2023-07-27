//go:build integration
// +build integration

package v1

import (
	"context"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/mehdihadeli/go-mediatr"
	"github.com/stretchr/testify/suite"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/test/messaging"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/test/messaging/consumer"
	testUtils "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/test/utils"
	dtosV1 "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orderservice/internal/orders/dtos/v1"
	createOrderCommandV1 "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orderservice/internal/orders/features/creating_order/v1/commands"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orderservice/internal/orders/features/creating_order/v1/dtos"
	integrationEvents "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orderservice/internal/orders/features/creating_order/v1/events/integration_events"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orderservice/internal/shared/test_fixtures/integration"
)

type createOrderIntegrationTests struct {
	*integration.IntegrationTestSharedFixture
}

func TestCreateOrderIntegration(t *testing.T) {
	suite.Run(
		t,
		&createOrderIntegrationTests{
			IntegrationTestSharedFixture: integration.NewIntegrationTestSharedFixture(t),
		},
	)
}

func (c *createOrderIntegrationTests) Test_Should_Create_New_Order_To_EventStoreDB() {
	command := createOrderCommandV1.NewCreateOrder(
		[]*dtosV1.ShopItemDto{
			{
				Quantity:    uint64(gofakeit.Number(1, 10)),
				Description: gofakeit.AdjectiveDescriptive(),
				Price:       gofakeit.Price(100, 10000),
				Title:       gofakeit.Name(),
			},
		},
		gofakeit.Email(),
		gofakeit.Address().Address,
		time.Now(),
	)
	result, err := mediatr.Send[*createOrderCommandV1.CreateOrder, *dtos.CreateOrderResponseDto](
		context.Background(),
		command,
	)

	c.NoError(err)
	c.NotNil(result)
	c.Equal(command.OrderId, result.OrderId)
}

func (c *createOrderIntegrationTests) Test_Should_Create_New_Order_To_MongoDB_Read() {
	ctx := context.Background()
	command := createOrderCommandV1.NewCreateOrder(
		[]*dtosV1.ShopItemDto{
			{
				Quantity:    uint64(gofakeit.Number(1, 10)),
				Description: gofakeit.AdjectiveDescriptive(),
				Price:       gofakeit.Price(100, 10000),
				Title:       gofakeit.Name(),
			},
		},
		gofakeit.Email(),
		gofakeit.Address().Address,
		time.Now(),
	)
	result, err := mediatr.Send[*createOrderCommandV1.CreateOrder, *dtos.CreateOrderResponseDto](
		context.Background(),
		command,
	)

	c.NoError(err)

	c.NoError(testUtils.WaitUntilConditionMet(func() bool {
		orderReadModel, err := c.OrderMongoRepository.GetOrderByOrderId(ctx, result.OrderId)
		c.NoError(err)
		return orderReadModel != nil
	}))
}

func (c *createOrderIntegrationTests) Test_Should_Publish_Order_Created_To_Broker() {
	ctx := context.Background()
	shouldPublish := messaging.ShouldProduced[*integrationEvents.OrderCreatedV1](
		ctx,
		c.Bus,
		nil,
	)

	command := createOrderCommandV1.NewCreateOrder(
		[]*dtosV1.ShopItemDto{
			{
				Quantity:    uint64(gofakeit.Number(1, 10)),
				Description: gofakeit.AdjectiveDescriptive(),
				Price:       gofakeit.Price(100, 10000),
				Title:       gofakeit.Name(),
			},
		},
		gofakeit.Email(),
		gofakeit.Address().Address,
		time.Now(),
	)
	_, err := mediatr.Send[*createOrderCommandV1.CreateOrder, *dtos.CreateOrderResponseDto](
		context.Background(),
		command,
	)
	c.NoError(err)

	// ensuring message published to the rabbitmq broker
	shouldPublish.Validate(ctx, "there is no published message", time.Second*30)
}

func (c *createOrderIntegrationTests) Test_Should_Consume_Order_Created_With_Existing_Consumer_From_Broker() {
	ctx := context.Background()

	// we setup this handler in `BeforeTest`
	// we don't have a consumer in this service, so we simulate one consumer
	// check for consuming `OrderCreatedV1` message with existing consumer
	hypothesis := messaging.ShouldConsume[*integrationEvents.OrderCreatedV1](ctx, c.Bus, nil)

	command := createOrderCommandV1.NewCreateOrder(
		[]*dtosV1.ShopItemDto{
			{
				Quantity:    uint64(gofakeit.Number(1, 10)),
				Description: gofakeit.AdjectiveDescriptive(),
				Price:       gofakeit.Price(100, 10000),
				Title:       gofakeit.Name(),
			},
		},
		gofakeit.Email(),
		gofakeit.Address().Address,
		time.Now(),
	)
	_, err := mediatr.Send[*createOrderCommandV1.CreateOrder, *dtos.CreateOrderResponseDto](
		context.Background(),
		command,
	)
	c.NoError(err)

	// ensuring message can be consumed with a consumer
	hypothesis.Validate(ctx, "there is no consumed message", time.Second*30)
}

func (c *createOrderIntegrationTests) Test_Should_Consume_Order_Created_With_New_Consumer_From_Broker() {
	ctx := context.Background()
	defer c.Bus.Stop()

	// check for consuming `OrderCreatedV1` message, with a new consumer
	hypothesis, err := messaging.ShouldConsumeNewConsumer[*integrationEvents.OrderCreatedV1](
		c.Bus,
	)
	c.Require().NoError(err)

	// at first, we should add new consumer to rabbitmq bus then start the broker, because we can't add new consumer after start.
	// we should also turn off consumer in `BeforeTest` for this test
	c.Bus.Start(ctx)

	// wait for consumers ready to consume before publishing messages, preparation background workers takes a bit time (for preventing messages lost)
	time.Sleep(1 * time.Second)

	command := createOrderCommandV1.NewCreateOrder(
		[]*dtosV1.ShopItemDto{
			{
				Quantity:    uint64(gofakeit.Number(1, 10)),
				Description: gofakeit.AdjectiveDescriptive(),
				Price:       gofakeit.Price(100, 10000),
				Title:       gofakeit.Name(),
			},
		},
		gofakeit.Email(),
		gofakeit.Address().Address,
		time.Now(),
	)
	_, err = mediatr.Send[*createOrderCommandV1.CreateOrder, *dtos.CreateOrderResponseDto](
		context.Background(),
		command,
	)
	c.NoError(err)

	// ensuring message can be consumed with a consumer
	hypothesis.Validate(ctx, "there is no consumed message", time.Second*30)
}

func (c *createOrderIntegrationTests) BeforeTest(suiteName, testName string) {
	if testName == "Test_Should_Consume_Order_Created_With_New_Consumer_From_Broker" {
		c.Bus.Stop()
	}
}

func (c *createOrderIntegrationTests) SetupSuite() {
	// we don't have a consumer in this service, so we simulate one consumer, register one consumer for `OrderCreatedV1` message before executing the tests
	testConsumer := consumer.NewRabbitMQFakeTestConsumerHandler[*integrationEvents.OrderCreatedV1]()
	err := c.Bus.ConnectConsumerHandler(&integrationEvents.OrderCreatedV1{}, testConsumer)
	c.Require().NoError(err)

	// in test mode we set rabbitmq `AutoStart=false`, so we should run rabbitmq bus manually
	c.Bus.Start(context.Background())
	// wait for consumers ready to consume before publishing messages, preparation background workers takes a bit time (for preventing messages lost)
	time.Sleep(1 * time.Second)
}

func (c *createOrderIntegrationTests) TearDownSuite() {
	c.Bus.Stop()
}
