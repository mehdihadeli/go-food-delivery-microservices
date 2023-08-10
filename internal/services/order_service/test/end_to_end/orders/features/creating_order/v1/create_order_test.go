//go:build e2e
// +build e2e

package v1

import (
	"net/http"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gavv/httpexpect/v2"
	customTypes "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/core/custom_types"
	"github.com/stretchr/testify/suite"

	dtosV1 "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orderservice/internal/orders/dtos/v1"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orderservice/internal/orders/features/creating_order/v1/dtos"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orderservice/internal/shared/test_fixtures/integration"
)

type createOrderE2ETest struct {
	*integration.IntegrationTestSharedFixture
}

func TestCreateProductE2E(t *testing.T) {
	suite.Run(
		t,
		&createOrderE2ETest{
			IntegrationTestSharedFixture: integration.NewIntegrationTestSharedFixture(t),
		},
	)
}

func (c *createOrderE2ETest) Test_Should_Return_Created_Status_With_Valid_Input() {
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

	// create httpexpect instance
	expect := httpexpect.Default(c.T(), c.BaseAddress)

	expect.POST("orders").
		WithJSON(request).
		Expect().
		Status(http.StatusCreated)
}
