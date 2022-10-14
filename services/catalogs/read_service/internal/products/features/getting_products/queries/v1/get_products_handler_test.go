package v1

import (
	"github.com/mehdihadeli/go-mediatr"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/test"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/utils"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/features/getting_products/dtos"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/shared/test_fixture/integration"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Get_Products_Query_Handler(t *testing.T) {
	test.SkipCI(t)
	fixture := integration.NewIntegrationTestFixture()

	err := mediatr.RegisterRequestHandler[*GetProducts, *dtos.GetProductsResponseDto](NewGetProductsHandler(fixture.Log, fixture.Cfg, fixture.MongoProductRepository))
	assert.NoError(t, err)

	fixture.Run()
	defer fixture.Cleanup()

	query := NewGetProducts(utils.NewListQuery(10, 1))
	queryResult, err := mediatr.Send[*GetProducts, *dtos.GetProductsResponseDto](fixture.Ctx, query)

	assert.NotNil(t, queryResult)
	assert.NotNil(t, queryResult.Products)
	assert.NotEmpty(t, queryResult.Products.Items)
}
