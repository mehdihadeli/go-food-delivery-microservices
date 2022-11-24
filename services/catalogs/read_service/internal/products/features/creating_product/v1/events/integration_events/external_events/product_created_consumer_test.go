package externalEvents

import (
	"testing"
	"time"

	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/features/creating_product/v1/commands"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/features/creating_product/v1/dtos"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/mehdihadeli/go-mediatr"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/types"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/test/messaging/consumer"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/test/utils"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/models"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/shared/test_fixture/integration"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func Test_Product_Created_Consumer_Should_Consume_Product_Created(t *testing.T) {
	testUtils.SkipCI(t)
	fixture := integration.NewIntegrationTestFixture()

	fakeConsumer := consumer.NewRabbitMQFakeTestConsumerHandler()
	err := fixture.Bus.ConnectConsumerHandler(ProductCreatedV1{}, fakeConsumer)
	assert.NoError(t, err)

	fixture.Run()
	defer fixture.Cleanup()

	err = fixture.Bus.PublishMessage(fixture.Ctx, &ProductCreatedV1{Message: types.NewMessage(uuid.NewV4().String())}, nil)
	assert.NoError(t, err)

	// ensuring message published to the rabbitmq broker
	assert.NoError(t, testUtils.WaitUntilConditionMet(func() bool {
		return fakeConsumer.IsHandled()
	}))
}

func Test_Product_Created_Consumer(t *testing.T) {
	testUtils.SkipCI(t)
	fixture := integration.NewIntegrationTestFixture()

	err := mediatr.RegisterRequestHandler[*commands.CreateProduct, *dtos.CreateProductResponseDto](commands.NewCreateProductHandler(fixture.Log, fixture.Cfg, fixture.MongoProductRepository, fixture.RedisProductRepository))
	assert.NoError(t, err)

	cons := NewProductCreatedConsumer(fixture.InfrastructureConfigurations)
	err = fixture.Bus.ConnectConsumerHandler(ProductCreatedV1{}, cons)
	assert.NoError(t, err)

	fixture.Run()
	defer fixture.Cleanup()

	pid := uuid.NewV4().String()
	productCreated := &ProductCreatedV1{
		Message:     types.NewMessage(uuid.NewV4().String()),
		ProductId:   pid,
		CreatedAt:   time.Now(),
		Name:        gofakeit.Name(),
		Price:       gofakeit.Price(150, 6000),
		Description: gofakeit.AdjectiveDescriptive(),
	}

	err = fixture.Bus.PublishMessage(fixture.Ctx, productCreated, nil)
	assert.NoError(t, err)

	var p *models.Product

	assert.NoError(t, testUtils.WaitUntilConditionMet(func() bool {
		p, err = fixture.MongoProductRepository.GetProductByProductId(fixture.Ctx, pid)
		assert.NoError(t, err)

		return p != nil
	}))

	assert.NotNil(t, p)
	assert.Equal(t, pid, p.ProductId)
}
