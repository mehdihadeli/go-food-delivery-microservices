package v1

import (
	"context"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/mehdihadeli/go-mediatr"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/test"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/shared/test_fixture/integration"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Update_Product_Command_Handler(t *testing.T) {
	test.SkipCI(t)
	fixture := integration.NewIntegrationTestFixture()

	err := mediatr.RegisterRequestHandler[*UpdateProduct, *mediatr.Unit](NewUpdateProductHandler(fixture.Log, fixture.Cfg, fixture.MongoProductRepository, fixture.RedisProductRepository))
	if err != nil {
		return
	}

	fixture.Run()
	defer fixture.Cleanup()

	productId, err := uuid.FromString("5f2c76c3-8f73-453c-af43-6b2a9551ff39")
	if err != nil {
		return
	}
	command := NewUpdateProduct(productId, gofakeit.Name(), gofakeit.AdjectiveDescriptive(), gofakeit.Price(150, 6000))
	result, err := mediatr.Send[*UpdateProduct, *mediatr.Unit](context.Background(), command)

	assert.NoError(t, err)
	assert.NotNil(t, result)
}
