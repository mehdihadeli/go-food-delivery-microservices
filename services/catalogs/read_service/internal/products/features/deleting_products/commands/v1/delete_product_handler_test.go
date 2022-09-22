package v1

import (
	"context"
	"github.com/mehdihadeli/go-mediatr"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/test"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/shared/test_fixture/integration"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Delete_Product_Command_Handler(t *testing.T) {
	test.SkipCI(t)
	fixture := integration.NewIntegrationTestFixture()

	err := mediatr.RegisterRequestHandler[*DeleteProduct, *mediatr.Unit](NewDeleteProductHandler(fixture.Log, fixture.Cfg, fixture.MongoProductRepository, fixture.RedisProductRepository))
	if err != nil {
		return
	}

	fixture.Run()
	defer fixture.Cleanup()

	productId, err := uuid.FromString("7f545256-4f20-4ef3-bdff-dd3c8e4a5408")
	if err != nil {
		return
	}
	command := NewDeleteProduct(productId)
	result, err := mediatr.Send[*DeleteProduct, *mediatr.Unit](context.Background(), command)

	assert.NoError(t, err)
	assert.NotNil(t, result)
}
