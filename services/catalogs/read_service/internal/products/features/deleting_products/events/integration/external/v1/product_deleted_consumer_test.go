package v1

import (
	"github.com/mehdihadeli/go-mediatr"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/types"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/test"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/test/messaging/consumer"
	v12 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/features/deleting_products/commands/v1"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/models"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/shared/test_fixture/integration"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Product_Deleted_Consumer_Should_Consume_Product_Deleted(t *testing.T) {
	test.SkipCI(t)
	fixture := integration.NewIntegrationTestFixture()

	fakeConsumer := consumer.NewRabbitMQFakeTestConsumerHandler()
	err := fixture.Bus.ConnectConsumerHandler(ProductDeletedV1{}, fakeConsumer)
	assert.NoError(t, err)

	fixture.Run()
	defer fixture.Cleanup()

	err = fixture.Bus.PublishMessage(fixture.Ctx, &ProductDeletedV1{Message: types.NewMessage(uuid.NewV4().String())}, nil)
	assert.NoError(t, err)

	// ensuring message published to the rabbitmq broker
	assert.NoError(t, test.WaitUntilConditionMet(func() bool {
		return fakeConsumer.IsHandled()
	}))
}

func Test_Product_Deleted_Consumer(t *testing.T) {
	test.SkipCI(t)
	fixture := integration.NewIntegrationTestFixture()

	err := mediatr.RegisterRequestHandler[*v12.DeleteProduct, *mediatr.Unit](v12.NewDeleteProductHandler(fixture.Log, fixture.Cfg, fixture.MongoProductRepository, fixture.RedisProductRepository))
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

	assert.NoError(t, test.WaitUntilConditionMet(func() bool {
		p, err = fixture.MongoProductRepository.GetProductByProductId(fixture.Ctx, pid)
		assert.NoError(t, err)

		return p == nil
	}))

	assert.Nil(t, p)
}
