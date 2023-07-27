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

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/shared/test_fixtures/integration"
)

type getProductByIdE2ETest struct {
	*integration.IntegrationTestSharedFixture
}

func TestGetProductByIdE2E(t *testing.T) {
	suite.Run(
		t,
		&getProductByIdE2ETest{
			IntegrationTestSharedFixture: integration.NewIntegrationTestSharedFixture(t),
		},
	)
}

func (c *getProductByIdE2ETest) Test_Should_Return_Ok_Status_With_Valid_Id() {
	ctx := context.Background()

	// create httpexpect instance
	expect := httpexpect.New(c.T(), c.BaseAddress)

	id := c.Items[0].ProductId

	expect.GET("products/{id}").
		WithPath("id", id.String()).
		WithContext(ctx).
		Expect().
		Status(http.StatusOK)
}

// Input validations
func (c *getProductByIdE2ETest) Test_Should_Return_NotFound_Status_With_Invalid_Id() {
	ctx := context.Background()

	// create httpexpect instance
	expect := httpexpect.New(c.T(), c.BaseAddress)

	id := uuid.NewV4()

	expect.GET("products/{id}").
		WithPath("id", id.String()).
		WithContext(ctx).
		Expect().
		Status(http.StatusNotFound)
}
