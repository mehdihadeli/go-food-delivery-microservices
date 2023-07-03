package v1

import (
	"context"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/mehdihadeli/go-mediatr"
	"github.com/stretchr/testify/suite"

	dtosV1 "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orderservice/internal/orders/dtos/v1"
	createOrderCommandV1 "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orderservice/internal/orders/features/creating_order/v1/commands"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orderservice/internal/orders/features/creating_order/v1/dtos"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orderservice/internal/shared/test_fixtures/integration"

	customTypes "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/core/custom_types"
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

func (c *createOrderIntegrationTests) Test_Should_Create_New_Order_To_DB() {
	//fakeConsumer := consumer.NewRabbitMQFakeTestConsumerHandler()
	//err = fixture.Bus.ConnectConsumerHandler(integrationEvents.OrderCreatedV1{}, fakeConsumer)
	//if err != nil {
	//	return
	//}
	//
	//if err != nil {
	//	return
	//}
	//
	//fixture.Run()

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

	command := createOrderCommandV1.NewCreateOrder(
		orderDto.ShopItems,
		orderDto.AccountEmail,
		orderDto.DeliveryAddress,
		time.Time(orderDto.DeliveryTime),
	)
	result, err := mediatr.Send[*createOrderCommandV1.CreateOrder, *dtos.CreateOrderResponseDto](
		context.Background(),
		command,
	)

	c.NoError(err)
	c.NotNil(result)
	c.Equal(command.OrderId, result.OrderId)

	//assert.NotNil(t, result)
	//assert.Equal(t, command.OrderId, result.OrderId)
	//time.Sleep(time.Second * 2)
	//
	//// ensuring message published to the rabbitmq broker
	//assert.NoError(t, testUtils.WaitUntilConditionMet(func() bool {
	//	return fakeConsumer.IsHandled()
	//}))
}
