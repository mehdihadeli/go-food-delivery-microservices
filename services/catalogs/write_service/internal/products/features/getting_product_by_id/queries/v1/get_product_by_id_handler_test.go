package v1

import (
	"context"
	"github.com/mehdihadeli/go-mediatr"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/test"
	gettingProductByIdDtos "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/getting_product_by_id/dtos"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/shared/test_fixtures/integration"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Get_Product_By_Id_Query_Handler(t *testing.T) {
	test.SkipCI(t)
	fixture := integration.NewIntegrationTestFixture()

	err := mediatr.RegisterRequestHandler[*GetProductById, *gettingProductByIdDtos.GetProductByIdResponseDto](NewGetProductByIdHandler(fixture.Log, fixture.Cfg, fixture.ProductRepository))
	if err != nil {
		return
	}

	fixture.Run()
	defer fixture.Cleanup()

	id, _ := uuid.FromString("1b088075-53f0-4376-a491-ca6fe3a7f8fa")
	query := NewGetProductById(id)
	queryResult, err := mediatr.Send[*GetProductById, *gettingProductByIdDtos.GetProductByIdResponseDto](context.Background(), query)

	assert.NotNil(t, queryResult)
	assert.NotNil(t, queryResult.Product)
	assert.Equal(t, id, queryResult.Product.ProductId)
}
