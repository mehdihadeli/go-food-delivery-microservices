package queries

import (
	"context"
	"testing"

	"github.com/mehdihadeli/go-mediatr"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"

	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/test/utils"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/features/getting_order_by_id/v1/dtos"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/shared/test_fixtures/integration"
)

func Test_Get_Order_By_Id_Query_Handler(t *testing.T) {
	testUtils.SkipCI(t)
	fixture := integration.NewIntegrationTestFixture()

	err := mediatr.RegisterRequestHandler[*GetOrderById, *dtos.GetOrderByIdResponseDto](NewGetOrderByIdHandler(fixture.Log(), fixture.Cfg(), fixture.MongoOrderReadRepository))
	if err != nil {
		return
	}

	fixture.Run()
	defer fixture.Cleanup()

	id, _ := uuid.FromString("22c09184-09b5-4dec-b70b-410b1d817ccc")
	query := NewGetOrderById(id)
	queryResult, err := mediatr.Send[*GetOrderById, *dtos.GetOrderByIdResponseDto](context.Background(), query)

	assert.NotNil(t, queryResult)
	assert.NotNil(t, queryResult.Order)
	assert.Equal(t, id.String(), queryResult.Order.Id)
	assert.NotNil(t, queryResult.Order.OrderId)
}
