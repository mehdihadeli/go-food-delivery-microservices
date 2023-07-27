//go:build e2e
// +build e2e

package v1

import (
	"context"
	"net/http"
	"testing"

	"github.com/gavv/httpexpect/v2"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/suite"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/mocks/testData"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/shared/test_fixtures/integration"
)

type deleteProductE2ETest struct {
	*integration.IntegrationTestSharedFixture
}

func TestDeleteProductE2E(t *testing.T) {
	suite.Run(
		t,
		&deleteProductE2ETest{
			IntegrationTestSharedFixture: integration.NewIntegrationTestSharedFixture(t),
		},
	)
}

func (c *deleteProductE2ETest) Test_Should_Return_No_Content_Status_With_Valid_Input() {
	ctx := context.Background()
	// create httpexpect instance
	expect := httpexpect.New(c.T(), c.BaseAddress)

	id := testData.Products[0].ProductId

	expect.DELETE("products/{id}").
		WithContext(ctx).
		WithPath("id", id.String()).
		Expect().
		Status(http.StatusNoContent)
}

// Input validations
func (c *deleteProductE2ETest) Test_Should_Return_Not_Found_Status_With_Invalid_Id() {
	ctx := context.Background()

	// create httpexpect instance
	expect := httpexpect.New(c.T(), c.BaseAddress)

	id := uuid.NewV4()

	expect.DELETE("products/{id}").
		WithContext(ctx).
		WithPath("id", id.String()).
		Expect().
		Status(http.StatusNotFound)
}
