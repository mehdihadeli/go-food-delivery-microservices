package v1

import (
	"context"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gavv/httpexpect/v2"
	customTypes "github.com/mehdihadeli/store-golang-microservice-sample/pkg/core/custom_types"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/test"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/delivery"
	ordersDto "github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/dtos"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/features/creating_order/dtos"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/shared/test_fixtures/e2e"
	"net/http"
	"testing"
	"time"
)

type My struct {
	Data string
	Num  int
}

// we could also run the server on docker and then send rest call to the api
func Test_Create_Order_E2E(t *testing.T) {
	test.SkipCI(t)

	fixture := e2e.NewE2ETestFixture()
	e := NewCreteOrderEndpoint(delivery.NewOrderEndpointBase(fixture.InfrastructureConfiguration, fixture.V1.OrdersGroup))
	e.MapRoute()

	defer fixture.Cleanup()

	fixture.Run()

	// create httpexpect instance
	expect := httpexpect.New(t, fixture.HttpServer.URL)

	request := dtos.CreateOrderRequestDto{
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

	expect.POST("/api/v1/orders").
		WithContext(context.Background()).
		WithJSON(request).
		Expect().
		Status(http.StatusCreated)

	time.Sleep(time.Second * 5)
}
