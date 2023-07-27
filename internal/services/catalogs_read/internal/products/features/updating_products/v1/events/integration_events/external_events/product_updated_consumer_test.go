//go:build.sh integration
// +build.sh integration

package externalEvents

import (
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/mehdihadeli/go-mediatr"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/messaging/types"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/test/messaging/consumer"
	testUtils "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/test/utils"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/read_service/internal/products/features/updating_products/v1/commands"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/read_service/internal/products/models"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/read_service/internal/shared/test_fixture/integration"
)

func Test_Product_Updated_Consumer_Should_Consume_Product_Updated(t *testing.T) {
	testUtils.SkipCI(t)
	fixture := integration.NewIntegrationTestFixture(integration.NewIntegrationTestSharedFixture(t))

	fakeConsumer := consumer.NewRabbitMQFakeTestConsumerHandler()
	err := fixture.Bus.ConnectConsumerHandler(ProductUpdatedV1{}, fakeConsumer)
	assert.NoError(t, err)

	fixture.Run()

	err = fixture.Bus.PublishMessage(
		fixture.Ctx,
		&ProductUpdatedV1{Message: types.NewMessage(uuid.NewV4().String())},
		nil,
	)
	assert.NoError(t, err)

	// ensuring message published to the rabbitmq broker
	assert.NoError(t, testUtils.WaitUntilConditionMet(func() bool {
		return fakeConsumer.IsHandled()
	}))
}

func Test_Product_Updated_Consumer(t *testing.T) {
	testUtils.SkipCI(t)
	fixture := integration.NewIntegrationTestFixture(integration.NewIntegrationTestSharedFixture(t))

	err := mediatr.RegisterRequestHandler[*commands.UpdateProduct, *mediatr.Unit](
		commands.NewUpdateProductHandler(
			fixture.Log,
			fixture.Cfg,
			fixture.MongoProductRepository,
			fixture.RedisProductRepository,
		),
	)
	assert.NoError(t, err)

	cons := NewProductUpdatedConsumer(fixture.InfrastructureConfigurations)
	err = fixture.Bus.ConnectConsumerHandler(ProductUpdatedV1{}, cons)
	assert.NoError(t, err)

	fixture.Run()

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

	assert.NoError(t, testUtils.WaitUntilConditionMet(func() bool {
		p, err = fixture.MongoProductRepository.GetProductByProductId(fixture.Ctx, pid)

		return p != nil && p.Name == productUpdated.Name
	}))

	assert.NotNil(t, p)
	assert.Equal(t, productUpdated.Name, p.Name)
	assert.Equal(t, productUpdated.ProductId, p.ProductId)
}
