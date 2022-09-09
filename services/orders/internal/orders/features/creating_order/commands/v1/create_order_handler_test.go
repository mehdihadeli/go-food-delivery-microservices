package v1

import (
	"context"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/mehdihadeli/go-mediatr"
	customTypes "github.com/mehdihadeli/store-golang-microservice-sample/pkg/core/custom_types"
	ordersDto "github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/dtos"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/features/creating_order/dtos"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/shared/test_fixtures/integration"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_Create_Order_Command_Handler(t *testing.T) {
	fixture := integration.NewIntegrationTestFixture()

	err := mediatr.RegisterRequestHandler[*CreateOrderCommand, *dtos.CreateOrderResponseDto](NewCreateOrderCommandHandler(fixture.Log, fixture.Cfg, fixture.OrderAggregateStore))
	if err != nil {
		return
	}

	defer fixture.Cleanup()

	orderDto := dtos.CreateOrderRequestDto{
		AccountEmail:    gofakeit.Email(),
		DeliveryAddress: gofakeit.Address().Address,
		DeliveryTime:    customTypes.CustomTime(time.Now()),
		ShopItems: []*ordersDto.ShopItemDto{
			{
				Quantity:    uint64(gofakeit.Number(1, 10)),
				Description: gofakeit.AdjectiveDescriptive(),
				Price:       gofakeit.Price(100, 10000),
				Title:       gofakeit.Name(),
			},
		},
	}

	command := NewCreateOrderCommand(orderDto.ShopItems, orderDto.AccountEmail, orderDto.DeliveryAddress, time.Time(orderDto.DeliveryTime))
	result, err := mediatr.Send[*CreateOrderCommand, *dtos.CreateOrderResponseDto](context.Background(), command)

	assert.NotNil(t, result)
	assert.Equal(t, command.OrderID, result.OrderID)
}
