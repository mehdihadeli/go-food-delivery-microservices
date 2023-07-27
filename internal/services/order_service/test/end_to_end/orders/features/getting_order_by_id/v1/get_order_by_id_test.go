package v1

import (
	"context"
	"net/http"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/suite"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orderservice/internal/shared/test_fixtures/integration"
)

type getOrderByIdE2ETest struct {
	*integration.IntegrationTestSharedFixture
}

func TestGetProductByIdE2E(t *testing.T) {
	suite.Run(
		t,
		&getOrderByIdE2ETest{
			IntegrationTestSharedFixture: integration.NewIntegrationTestSharedFixture(t),
		},
	)
}

func (c *getOrderByIdE2ETest) Test_Should_Return_Ok_Status_With_Valid_Id() {
	ctx := context.Background()

	expect := httpexpect.Default(c.T(), c.BaseAddress)

	id := c.Items[0].Id

	expect.GET("orders/{id}").
		WithPath("id", id).
		WithContext(ctx).
		Expect().
		Status(http.StatusOK)
}
