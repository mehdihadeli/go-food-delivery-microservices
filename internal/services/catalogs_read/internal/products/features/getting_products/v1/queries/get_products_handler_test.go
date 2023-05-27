package queries

import (
	"testing"

	"github.com/mehdihadeli/go-mediatr"
	"github.com/stretchr/testify/assert"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/test/utils"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/utils"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/read_service/internal/products/features/getting_products/v1/dtos"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/read_service/internal/shared/test_fixture/integration"
)

func Test_Get_Products_Query_Handler(t *testing.T) {
	testUtils.SkipCI(t)
	fixture := integration.NewIntegrationTestFixture(integration.NewIntegrationTestSharedFixture(t))

	err := mediatr.RegisterRequestHandler[*GetProducts, *dtos.GetProductsResponseDto](NewGetProductsHandler(fixture.Log, fixture.Cfg, fixture.MongoProductRepository))
	assert.NoError(t, err)

	fixture.Run()

	query := NewGetProducts(utils.NewListQuery(10, 1))
	queryResult, err := mediatr.Send[*GetProducts, *dtos.GetProductsResponseDto](fixture.Ctx, query)

	assert.NotNil(t, queryResult)
	assert.NotNil(t, queryResult.Products)
	assert.NotEmpty(t, queryResult.Products.Items)
}
