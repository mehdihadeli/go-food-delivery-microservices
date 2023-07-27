//go:build.sh integration
// +build.sh integration

package queries

import (
	"testing"

	"github.com/mehdihadeli/go-mediatr"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"

	testUtils "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/test/utils"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/read_service/internal/products/features/get_product_by_id/v1/dtos"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/read_service/internal/shared/test_fixture/integration"
)

func Test_Get_Product_By_Id_Query_Handler(t *testing.T) {
	testUtils.SkipCI(t)
	fixture := integration.NewIntegrationTestFixture(integration.NewIntegrationTestSharedFixture(t))

	err := mediatr.RegisterRequestHandler[*GetProductById, *dtos.GetProductByIdResponseDto](
		NewGetProductByIdHandler(
			fixture.Log,
			fixture.Cfg,
			fixture.MongoProductRepository,
			fixture.RedisProductRepository,
		),
	)
	assert.NoError(t, err)

	fixture.Run()

	id, err := uuid.FromString("86093212-2e4c-4316-b1ef-f545154ba40d")
	assert.NoError(t, err)

	query := NewGetProductById(id)
	result, err := mediatr.Send[*GetProductById, *dtos.GetProductByIdResponseDto](
		fixture.Ctx,
		query,
	)

	assert.NotNil(t, result.Product)
	assert.Equal(t, result.Product.Id, id.String())
}
