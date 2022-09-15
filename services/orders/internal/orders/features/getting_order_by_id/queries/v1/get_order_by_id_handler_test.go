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

	err := mediatr.RegisterRequestHandler[*GetOrderByIdQuery, *dtos.GetOrderByIdResponseDto](NewGetOrderByIdHandler(fixture.Log, fixture.Cfg, fixture.MongoOrderReadRepository))
	if err != nil {
		return
	}

	defer fixture.Cleanup()

	id, _ := uuid.FromString("1b4b0599-bc3c-4c1d-94af-fd1895713620")
	query := NewGetOrderByIdQuery(id)
	queryResult, err := mediatr.Send[*GetOrderByIdQuery, *dtos.GetOrderByIdResponseDto](context.Background(), query)

	assert.NotNil(t, queryResult)
	assert.NotNil(t, queryResult.Order)
	assert.Equal(t, id.String(), queryResult.Order.Id)
}
