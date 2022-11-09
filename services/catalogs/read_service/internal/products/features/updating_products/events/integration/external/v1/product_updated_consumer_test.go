package v1

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/mehdihadeli/go-mediatr"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/types"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/test"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/test/messaging/consumer"
	v1 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/features/updating_products/commands/v1"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/models"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/shared/test_fixture/integration"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_Product_Updated_Consumer_Should_Consume_Product_Updated(t *testing.T) {
	test.SkipCI(t)
	fixture := integration.NewIntegrationTestFixture()

	fakeConsumer := consumer.NewRabbitMQFakeTestConsumerHandlerWithHypothesis()
	err := fixture.Bus.ConnectConsumerHandler(ProductUpdatedV1{}, fakeConsumer)
	assert.NoError(t, err)

	fixture.Run()
	defer fixture.Cleanup()

	err = fixture.Bus.PublishMessage(fixture.Ctx, &ProductUpdatedV1{Message: types.NewMessage(uuid.NewV4().String())}, nil)
	assert.NoError(t, err)

	// ensuring message published to the rabbitmq broker
	assert.NoError(t, test.WaitUntilConditionMet(func() bool {
		return fakeConsumer.IsHandled()
	}))
}

func Test_Product_Updated_Consumer(t *testing.T) {
	test.SkipCI(t)
	fixture := integration.NewIntegrationTestFixture()

	err := mediatr.RegisterRequestHandler[*v1.UpdateProduct, *mediatr.Unit](v1.NewUpdateProductHandler(fixture.Log, fixture.Cfg, fixture.MongoProductRepository, fixture.RedisProductRepository))
	assert.NoError(t, err)

	cons := NewProductUpdatedConsumer(fixture.InfrastructureConfigurations)
	err = fixture.Bus.ConnectConsumerHandler(ProductUpdatedV1{}, cons)
	assert.NoError(t, err)

	fixture.Run()
	defer fixture.Cleanup()

	pid := "74dde762-7676-43e9-a849-3b6317cc722b"

	productUpdated := &ProductUpdatedV1{
		Message:     types.NewMessage(uuid.NewV4().String()),
		ProductId:   pid,
		Name:        gofakeit.Name(),
		Description: gofakeit.AdjectiveDescriptive(),
		Price:       gofakeit.Price(150, 6000),
		UpdatedAt:   time.Now(),
	}

	err = fixture.Bus.PublishMessage(fixture.Ctx, productUpdated, nil)
	assert.NoError(t, err)

	var p *models.Product

	assert.NoError(t, test.WaitUntilConditionMet(func() bool {
		p, err = fixture.MongoProductRepository.GetProductByProductId(fixture.Ctx, pid)

		return p != nil && p.Name == productUpdated.Name
	}))

	assert.NotNil(t, p)
	assert.Equal(t, productUpdated.Name, p.Name)
	assert.Equal(t, productUpdated.ProductId, p.ProductId)
}
