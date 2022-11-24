package commands

import (
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/mehdihadeli/go-mediatr"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"

	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/test/utils"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/shared/test_fixture/integration"
)

func Test_Update_Product_Command_Handler(t *testing.T) {
	testUtils.SkipCI(t)
	fixture := integration.NewIntegrationTestFixture()

	err := mediatr.RegisterRequestHandler[*UpdateProduct, *mediatr.Unit](NewUpdateProductHandler(fixture.Log, fixture.Cfg, fixture.MongoProductRepository, fixture.RedisProductRepository))
	assert.NoError(t, err)

	fixture.Run()
	defer fixture.Cleanup()

	productId, err := uuid.FromString("34dac034-ad17-427d-9bc1-3d7dc07c40f0")
	assert.NoError(t, err)

	command := NewUpdateProduct(productId, gofakeit.Name(), gofakeit.AdjectiveDescriptive(), gofakeit.Price(150, 6000))
	result, err := mediatr.Send[*UpdateProduct, *mediatr.Unit](fixture.Ctx, command)
	assert.NoError(t, err)

	assert.NotNil(t, result)
}
