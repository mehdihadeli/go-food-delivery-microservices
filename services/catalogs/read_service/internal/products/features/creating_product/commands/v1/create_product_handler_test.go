package v1

import (
	"context"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/mehdihadeli/go-mediatr"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/test"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/features/creating_product"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/shared/test_fixture/integration"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_Create_Product_Command_Handler(t *testing.T) {
	test.SkipCI(t)
	fixture := integration.NewIntegrationTestFixture()

	err := mediatr.RegisterRequestHandler[*CreateProduct, *creating_product.CreateProductResponseDto](NewCreateProductHandler(fixture.Log, fixture.Cfg, fixture.MongoProductRepository, fixture.RedisProductRepository))
	if err != nil {
		return
	}

	fixture.Run()
	defer fixture.Cleanup()

	command := NewCreateProduct(gofakeit.UUID(), gofakeit.Name(), gofakeit.AdjectiveDescriptive(), gofakeit.Price(150, 6000), time.Now())
	result, err := mediatr.Send[*CreateProduct, *creating_product.CreateProductResponseDto](context.Background(), command)

	assert.NotNil(t, result)
	assert.NotEmpty(t, result.Id)
}
