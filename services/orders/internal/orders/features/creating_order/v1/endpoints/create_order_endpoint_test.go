package createOrderV1

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gavv/httpexpect/v2"

	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/core/custom_types"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/test/utils"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/delivery"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/dtos/v1"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/features/creating_order/v1/dtos"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/shared/test_fixtures/e2e"
)

func Test_Create_Order_E2E(t *testing.T) {
	testUtils.SkipCI(t)

	fixture := e2e.NewE2ETestFixture()
	e := NewCreteOrderEndpoint(delivery.NewOrderEndpointBase(fixture.InfrastructureConfigurations, fixture.V1.OrdersGroup, fixture.Bus, fixture.OrdersMetrics))
	e.MapRoute()

	fixture.Run()
	defer fixture.Cleanup()

	// create httpexpect instance
	expect := httpexpect.New(t, fixture.HttpServer.URL)

	request := dtos.CreateOrderRequestDto{
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

	expect.POST("/api/v1/orders").
		WithContext(context.Background()).
		WithJSON(request).
		Expect().
		Status(http.StatusCreated)
}
