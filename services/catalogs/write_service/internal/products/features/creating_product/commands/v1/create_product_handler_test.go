package v1

import (
	"context"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/mehdihadeli/go-mediatr"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/test"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/test/messaging/consumer"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/creating_product/dtos"
	v1 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/creating_product/events/integration/v1"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/shared/test_fixtures/integration"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Create_Product_Command_Handler(t *testing.T) {
	test.SkipCI(t)
	fixture := integration.NewIntegrationTestFixture()

	err := mediatr.RegisterRequestHandler[*CreateProduct, *dtos.CreateProductResponseDto](NewCreateProductHandler(fixture.Log(), fixture.Cfg(), fixture.ProductRepository, fixture.Bus))
	if err != nil {
		return
	}

	fakeConsumer := consumer.NewRabbitMQFakeTestConsumer()
	err = fixture.Bus.ConnectConsumer(v1.ProductCreatedV1{}, fakeConsumer)
	if err != nil {
		return
	}

	fixture.Run()
	defer fixture.Cleanup()

	command := NewCreateProduct(gofakeit.Name(), gofakeit.AdjectiveDescriptive(), gofakeit.Price(150, 6000))
	result, err := mediatr.Send[*CreateProduct, *dtos.CreateProductResponseDto](context.Background(), command)

	assert.NotNil(t, result)
	assert.Equal(t, command.ProductID, result.ProductID)

	// ensuring message published to the rabbitmq broker
	assert.NoError(t, test.WaitUntilConditionMet(func() bool {
		return fakeConsumer.IsHandled()
	}))
}
