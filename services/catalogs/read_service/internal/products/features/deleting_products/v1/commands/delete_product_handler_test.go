package commands

import (
	"context"
	"testing"

	"github.com/mehdihadeli/go-mediatr"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"

	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/test/utils"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/shared/test_fixture/integration"
)

func Test_Delete_Product_Command_Handler(t *testing.T) {
	testUtils.SkipCI(t)
	fixture := integration.NewIntegrationTestFixture()

	err := mediatr.RegisterRequestHandler[*DeleteProduct, *mediatr.Unit](NewDeleteProductHandler(fixture.Log, fixture.Cfg, fixture.MongoProductRepository, fixture.RedisProductRepository))
	assert.NoError(t, err)

	fixture.Run()
	defer fixture.Cleanup()

	productId, err := uuid.FromString("399beedb-0f2c-4dc6-b53c-51aa0a2f7a91")
	assert.NoError(t, err)

	command := NewDeleteProduct(productId)
	result, err := mediatr.Send[*DeleteProduct, *mediatr.Unit](context.Background(), command)

	assert.NoError(t, err)
	assert.NotNil(t, result)
}
