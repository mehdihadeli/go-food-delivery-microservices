//go:build e2e
// +build e2e

package v1

import (
	"context"
	"net/http"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/suite"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/features/creating_product/v1/dtos"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/shared/test_fixtures/integration"
)

type createProductE2ETest struct {
	*integration.IntegrationTestSharedFixture
}

func TestCreateProductE2E(t *testing.T) {
	suite.Run(
		t,
		&createProductE2ETest{
			IntegrationTestSharedFixture: integration.NewIntegrationTestSharedFixture(t),
		},
	)
}

func (c *createProductE2ETest) Test_Should_Return_Created_Status_With_Valid_Input() {
	request := dtos.CreateProductRequestDto{
		Description: gofakeit.AdjectiveDescriptive(),
		Price:       gofakeit.Price(100, 1000),
		Name:        gofakeit.Name(),
	}

	// create httpexpect instance
	expect := httpexpect.New(c.T(), c.BaseAddress)

	expect.POST("products").
		// WithContext(context.Background()).
		WithJSON(request).
		Expect().
		Status(http.StatusCreated)
}

// Input validations
func (c *createProductE2ETest) Test_Should_Return_Bad_Request_Status_With_Invalid_Price_Input() {
	request := dtos.CreateProductRequestDto{
		Description: gofakeit.AdjectiveDescriptive(),
		Price:       0,
		Name:        gofakeit.Name(),
	}

	// create httpexpect instance
	expect := httpexpect.New(c.T(), c.BaseAddress)

	expect.POST("products").
		WithContext(context.Background()).
		WithJSON(request).
		Expect().
		Status(http.StatusBadRequest)
}
