//go:build e2e
// +build e2e

package v1

import (
	"context"
	"net/http"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/suite"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/shared/test_fixtures/integration"
)

type getProductsE2ETest struct {
	*integration.IntegrationTestSharedFixture
}

func TestGetProductsE2e(t *testing.T) {
	suite.Run(
		t,
		&getProductsE2ETest{
			IntegrationTestSharedFixture: integration.NewIntegrationTestSharedFixture(t),
		},
	)
}

func (c *getProductsE2ETest) Test_Should_Return_Ok_Status() {
	ctx := context.Background()

	// create httpexpect instance
	expect := httpexpect.New(c.T(), c.BaseAddress)

	expect.GET("products").
		WithContext(ctx).
		Expect().
		Status(http.StatusOK)
}
