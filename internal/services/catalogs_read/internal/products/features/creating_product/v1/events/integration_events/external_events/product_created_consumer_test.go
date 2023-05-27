package externalEvents

import (
	"testing"
	"time"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/messaging/types"
	testUtils "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/test/utils"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/read_service/internal/products/features/creating_product/v1/commands"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/read_service/internal/products/features/creating_product/v1/dtos"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/mehdihadeli/go-mediatr"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/test/messaging/consumer"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/read_service/internal/products/models"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/read_service/internal/shared/test_fixture/integration"
)

func Test_Product_Created_Consumer_Should_Consume_Product_Created(t *testing.T) {
	testUtils.SkipCI(t)
	fixture := integration.NewIntegrationTestFixture(integration.NewIntegrationTestSharedFixture(t))

	fakeConsumer := consumer.NewRabbitMQFakeTestConsumerHandler()
	err := fixture.Bus.ConnectConsumerHandler(ProductCreatedV1{}, fakeConsumer)
	assert.NoError(t, err)

	fixture.Run()

	err = fixture.Bus.PublishMessage(fixture.Ctx, &ProductCreatedV1{Message: types.NewMessage(uuid.NewV4().String())}, nil)
	assert.NoError(t, err)

	// ensuring message published to the rabbitmq broker
	assert.NoError(t, testUtils.WaitUntilConditionMet(func() bool {
		return fakeConsumer.IsHandled()
	}))
}

func Test_Product_Created_Consumer(t *testing.T) {
	testUtils.SkipCI(t)
	fixture := integration.NewIntegrationTestFixture(integration.NewIntegrationTestSharedFixture(t))

	err := mediatr.RegisterRequestHandler[*commands.CreateProduct, *dtos.CreateProductResponseDto](commands.NewCreateProductHandler(fixture.Log, fixture.Cfg, fixture.MongoProductRepository, fixture.RedisProductRepository))
	assert.NoError(t, err)

	cons := NewProductCreatedConsumer(fixture.InfrastructureConfigurations)
	err = fixture.Bus.ConnectConsumerHandler(ProductCreatedV1{}, cons)
	assert.NoError(t, err)

	fixture.Run()

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
