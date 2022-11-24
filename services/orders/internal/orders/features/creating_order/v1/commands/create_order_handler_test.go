package createOrderCommandV1

import (
	"context"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/mehdihadeli/go-mediatr"
	"github.com/stretchr/testify/assert"

	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/core/custom_types"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/test/messaging/consumer"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/test/utils"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/dtos/v1"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/features/creating_order/v1/dtos"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/features/creating_order/v1/events/integration_events"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/shared/test_fixtures/integration"
)

func Test_Create_Order_Command_Handler(t *testing.T) {
	testUtils.SkipCI(t)
	fixture := integration.NewIntegrationTestFixture()
	defer fixture.Cleanup()

	err := mediatr.RegisterRequestHandler[*CreateOrder, *dtos.CreateOrderResponseDto](NewCreateOrderHandler(fixture.Log(), fixture.Cfg(), fixture.OrderAggregateStore))
	if err != nil {
		return
	}

	fakeConsumer := consumer.NewRabbitMQFakeTestConsumerHandler()
	err = fixture.Bus.ConnectConsumerHandler(integrationEvents.OrderCreatedV1{}, fakeConsumer)
	if err != nil {
		return
	}

	if err != nil {
		return
	}

	fixture.Run()

	orderDto := dtos.CreateOrderRequestDto{
		AccountEmail:    gofakeit.Email(),
		DeliveryAddress: gofakeit.Address().Address,
		DeliveryTime:    customTypes.CustomTime(time.Now()),
		ShopItems: []*dtosV1.ShopItemDto{
			{
				Quantity:    uint64(gofakeit.Number(1, 10)),
				Description: gofakeit.AdjectiveDescriptive(),
				Price:       gofakeit.Price(100, 10000),
				Title:       gofakeit.Name(),
			},
		},
	}

	command := NewCreateOrder(orderDto.ShopItems, orderDto.AccountEmail, orderDto.DeliveryAddress, time.Time(orderDto.DeliveryTime))
	result, err := mediatr.Send[*CreateOrder, *dtos.CreateOrderResponseDto](context.Background(), command)

	assert.NotNil(t, result)
	assert.Equal(t, command.OrderId, result.OrderId)
	time.Sleep(time.Second * 2)

	// ensuring message published to the rabbitmq broker
	assert.NoError(t, testUtils.WaitUntilConditionMet(func() bool {
		return fakeConsumer.IsHandled()
	}))
}
