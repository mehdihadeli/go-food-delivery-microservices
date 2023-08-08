//go:build integration
// +build integration

package v1

import (
	"context"
	"testing"

	"github.com/mehdihadeli/go-mediatr"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/suite"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orderservice/internal/orders/features/getting_order_by_id/v1/dtos"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orderservice/internal/orders/features/getting_order_by_id/v1/queries"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orderservice/internal/shared/test_fixtures/integration"
)

type getOrderByIdIntegrationTests struct {
	*integration.IntegrationTestSharedFixture
}

func TestCreateOrderIntegration(t *testing.T) {
	suite.Run(
		t,
		&getOrderByIdIntegrationTests{
			IntegrationTestSharedFixture: integration.NewIntegrationTestSharedFixture(t),
		},
	)
}

func (c *getOrderByIdIntegrationTests) Test_Should_Returns_Existing_Order_From_DB_With_Correct_Properties() {
	ctx := context.Background()

	id, err := uuid.FromString(c.Items[0].Id)
	c.Require().NoError(err)

	query := queries.NewGetOrderById(id)

	result, err := mediatr.Send[*queries.GetOrderById, *dtos.GetOrderByIdResponseDto](
		ctx,
		query,
	)

	c.Require().NoError(err)
	c.NotNil(result)
	c.NotNil(result.Order)
	c.Equal(id.String(), result.Order.Id)
	c.NotNil(result.Order.OrderId)
}
