//go:build e2e
// +build e2e

package v1

import (
	"context"
	"net/http"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/suite"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/mocks/testData"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/shared/test_fixtures/integration"
)

type searchProductsE2ETest struct {
	*integration.IntegrationTestSharedFixture
}

func TestCreateProductE2e(t *testing.T) {
	suite.Run(
		t,
		&searchProductsE2ETest{
			IntegrationTestSharedFixture: integration.NewIntegrationTestSharedFixture(t),
		},
	)
}

func (c *searchProductsE2ETest) Test_Should_Return_Ok_Status() {
	ctx := context.Background()
	// create httpexpect instance
	expect := httpexpect.New(c.T(), c.BaseAddress)

	expect.GET("products/search").
		WithContext(ctx).
		WithQuery("search", testData.Products[0].Name).
		Expect().
		Status(http.StatusOK)
}
