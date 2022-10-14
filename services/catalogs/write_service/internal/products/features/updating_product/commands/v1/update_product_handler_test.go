package v1

import (
	"context"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/mehdihadeli/go-mediatr"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/test"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/test/messaging/consumer"
	v12 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/updating_product/events/integration/v1"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/shared/test_fixtures/integration"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Update_Product_Command_Handler(t *testing.T) {
	test.SkipCI(t)
	fixture := integration.NewIntegrationTestFixture()

	err := mediatr.RegisterRequestHandler[*UpdateProduct, *mediatr.Unit](NewUpdateProductHandler(fixture.Log, fixture.Cfg, fixture.ProductRepository, fixture.Bus))
	assert.NoError(t, err)

	fakeConsumer := consumer.NewRabbitMQFakeTestConsumerHandler()
	err = fixture.Bus.ConnectConsumerHandler(v12.ProductUpdatedV1{}, fakeConsumer)

	fixture.Run()
	defer fixture.Cleanup()

	id, err := uuid.FromString("49a8e487-945b-4050-9a4c-a9242247cb48")
	assert.NoError(t, err)

	command := NewUpdateProduct(id, gofakeit.Name(), gofakeit.AdjectiveDescriptive(), gofakeit.Price(150, 6000))
	result, err := mediatr.Send[*UpdateProduct, *mediatr.Unit](context.Background(), command)

	assert.NoError(t, err)
	assert.NotNil(t, result)

	// ensuring message published to the rabbitmq broker
	assert.NoError(t, test.WaitUntilConditionMet(func() bool {
		return fakeConsumer.IsHandled()
	}))
}
