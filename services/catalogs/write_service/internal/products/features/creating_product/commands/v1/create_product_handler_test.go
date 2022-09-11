package v1

import (
	"context"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/mehdihadeli/go-mediatr"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/test"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/creating_product/dtos"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/shared/test_fixtures/integration"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Create_Product_Command_Handler(t *testing.T) {
	test.SkipCI(t)
	fixture := integration.NewIntegrationTestFixture()

	err := mediatr.RegisterRequestHandler[*CreateProductCommand, *dtos.CreateProductResponseDto](NewCreateProductCommandHandler(fixture.Log, fixture.Cfg, fixture.ProductRepository, fixture.KafkaProducer))
	if err != nil {
		return
	}

	defer fixture.Cleanup()

	command := NewCreateProductCommand(gofakeit.Name(), gofakeit.AdjectiveDescriptive(), gofakeit.Price(150, 6000))
	result, err := mediatr.Send[*CreateProductCommand, *dtos.CreateProductResponseDto](context.Background(), command)

	assert.NotNil(t, result)
	assert.Equal(t, command.ProductID, result.ProductID)
}
