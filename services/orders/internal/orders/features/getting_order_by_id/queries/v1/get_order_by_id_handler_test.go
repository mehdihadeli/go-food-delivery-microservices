package v1

import (
	"context"
	"github.com/mehdihadeli/go-mediatr"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/test"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/features/getting_order_by_id/dtos"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/shared/test_fixtures/integration"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Get_Order_By_Id_Query_Handler(t *testing.T) {
	test.SkipCI(t)
	fixture := integration.NewIntegrationTestFixture()

	err := mediatr.RegisterRequestHandler[*GetOrderById, *dtos.GetOrderByIdResponseDto](NewGetOrderByIdHandler(fixture.Log, fixture.Cfg, fixture.MongoOrderReadRepository))
	if err != nil {
		return
	}

	fixture.Run()
	defer fixture.Cleanup()

	id, _ := uuid.FromString("c8018f1e-787b-4d5e-98fd-4b4e072d56b2")
	query := NewGetOrderById(id)
	queryResult, err := mediatr.Send[*GetOrderById, *dtos.GetOrderByIdResponseDto](context.Background(), query)

	assert.NotNil(t, queryResult)
	assert.NotNil(t, queryResult.Order)
	assert.Equal(t, id.String(), queryResult.Order.Id)
	assert.NotNil(t, queryResult.Order.OrderId)
}
