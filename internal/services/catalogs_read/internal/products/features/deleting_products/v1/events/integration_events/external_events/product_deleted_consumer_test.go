package externalEvents

import (
	"testing"

	"github.com/mehdihadeli/go-mediatr"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/messaging/types"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/test/messaging/consumer"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/read_service/internal/products/features/deleting_products/v1/commands"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/read_service/internal/products/models"
)

func Test_Product_Deleted_Consumer_Should_Consume_Product_Deleted(t *testing.T) {
	testUtils.SkipCI(t)
	fixture := integration.NewIntegrationTestFixture()

	fakeConsumer := consumer.NewRabbitMQFakeTestConsumerHandler()
	err := fixture.Bus.ConnectConsumerHandler(ProductDeletedV1{}, fakeConsumer)
	assert.NoError(t, err)

	fixture.Run()
	defer fixture.Cleanup()

	err = fixture.Bus.PublishMessage(fixture.Ctx, &ProductDeletedV1{Message: types.NewMessage(uuid.NewV4().String())}, nil)
	assert.NoError(t, err)

	// ensuring message published to the rabbitmq broker
	assert.NoError(t, testUtils.WaitUntilConditionMet(func() bool {
		return fakeConsumer.IsHandled()
	}))
}

func Test_Product_Deleted_Consumer(t *testing.T) {
	testUtils.SkipCI(t)
	fixture := integration.NewIntegrationTestFixture()

	err := mediatr.RegisterRequestHandler[*commands.DeleteProduct, *mediatr.Unit](commands.NewDeleteProductHandler(fixture.Log, fixture.Cfg, fixture.MongoProductRepository, fixture.RedisProductRepository))
	assert.NoError(t, err)

	cons := NewProductDeletedConsumer(fixture.InfrastructureConfigurations)
	err = fixture.Bus.ConnectConsumerHandler(ProductDeletedV1{}, cons)
	assert.NoError(t, err)

	fixture.Run()
	defer fixture.Cleanup()

	pid := "ff13c422-e0dc-466d-9bee-d09c1d3122e1"
	productDeleted := &ProductDeletedV1{
		Message:   types.NewMessage(uuid.NewV4().String()),
		ProductId: pid,
	}

	err = fixture.Bus.PublishMessage(fixture.Ctx, productDeleted, nil)
	assert.NoError(t, err)

	var p *models.Product

	assert.NoError(t, testUtils.WaitUntilConditionMet(func() bool {
		p, err = fixture.MongoProductRepository.GetProductByProductId(fixture.Ctx, pid)
		assert.NoError(t, err)

		return p == nil
	}))

	assert.Nil(t, p)
}
