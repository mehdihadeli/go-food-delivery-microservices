package v1

import (
	"context"
	"github.com/mehdihadeli/go-mediatr"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/features/getting_order_by_id/dtos"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/shared/test_fixtures/integration"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Get_Order_By_Id_Query_Handler(t *testing.T) {
	fixture := integration.NewIntegrationTestFixture()

	err := mediatr.RegisterRequestHandler[*GetOrderByIdQuery, *dtos.GetOrderByIdResponseDto](NewGetOrderByIdHandler(fixture.Log, fixture.Cfg, fixture.OrderAggregateStore))
	if err != nil {
		return
	}

	defer fixture.Cleanup()

	id, _ := uuid.FromString("97e2d953-ed25-4afb-8578-782cc5d365ba")
	query := NewGetOrderByIdQuery(id)
	queryResult, err := mediatr.Send[*GetOrderByIdQuery, *dtos.GetOrderByIdResponseDto](context.Background(), query)

	assert.NotNil(t, queryResult)
	assert.NotNil(t, queryResult.Order)
	assert.Equal(t, id, queryResult.Order.Id)
}
