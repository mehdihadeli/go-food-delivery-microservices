//go:build e2e
// +build e2e

package v1

import (
	"context"
	"net/http"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/suite"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogreadservice/internal/shared/test_fixture/integration"
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

	expect := httpexpect.New(c.T(), c.BaseAddress)

	id := c.Items[0].Id

	expect.GET("products/{id}").
		WithPath("id", id).
		WithContext(ctx).
		Expect().
		Status(http.StatusOK)
}
