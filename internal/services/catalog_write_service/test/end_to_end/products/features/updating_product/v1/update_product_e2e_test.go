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

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/features/updating_product/v1/dtos"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/shared/test_fixtures/integration"
)

type updateProductE2ETest struct {
	*integration.IntegrationTestSharedFixture
}

func TestCreateProductE2E(t *testing.T) {
	suite.Run(
		t,
		&updateProductE2ETest{
			IntegrationTestSharedFixture: integration.NewIntegrationTestSharedFixture(t),
		},
	)
}

func (c *updateProductE2ETest) Test_Should_Return_NoContent_Status_With_Valid_Input() {
	id := c.Items[0].ProductId

	request := dtos.UpdateProductRequestDto{
		Description: gofakeit.AdjectiveDescriptive(),
		Price:       gofakeit.Price(100, 1000),
		Name:        gofakeit.Name(),
	}

	// create httpexpect instance
	expect := httpexpect.New(c.T(), c.BaseAddress)

	expect.PUT("products/{id}").
		WithContext(context.Background()).
		WithJSON(request).
		WithPath("id", id.String()).
		Expect().
		Status(http.StatusNoContent)
}

// Input validations
func (c *updateProductE2ETest) Test_Should_Return_Bad_Request_Status_With_Invalid_Input() {
	id := c.Items[0].ProductId

	request := dtos.UpdateProductRequestDto{
		Description: gofakeit.AdjectiveDescriptive(),
		Price:       0,
		Name:        gofakeit.Name(),
	}

	// create httpexpect instance
	expect := httpexpect.New(c.T(), c.BaseAddress)

	expect.PUT("products/{id}").
		WithContext(context.Background()).
		WithJSON(request).
		WithPath("id", id.String()).
		Expect().
		Status(http.StatusBadRequest)
}
