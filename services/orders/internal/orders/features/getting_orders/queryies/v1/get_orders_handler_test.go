package v1

import (
	"context"
	"github.com/mehdihadeli/go-mediatr"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/test"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/utils"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/features/getting_orders/dtos"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/shared/test_fixtures/integration"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Get_Orders_Query_Handler(t *testing.T) {
	test.SkipCI(t)
	fixture := integration.NewIntegrationTestFixture()

	err := mediatr.RegisterRequestHandler[*GetOrders, *dtos.GetOrdersResponseDto](NewGetOrdersHandler(fixture.Log, fixture.Cfg, fixture.MongoOrderReadRepository))
	if err != nil {
		return
	}

	defer fixture.Cleanup()

	query := NewGetOrders(utils.NewListQuery(10, 1))
	queryResult, err := mediatr.Send[*GetOrders, *dtos.GetOrdersResponseDto](context.Background(), query)

	assert.NotNil(t, queryResult)
	assert.NotNil(t, queryResult.Orders)
	assert.NotEmpty(t, queryResult.Orders.Items)
}
