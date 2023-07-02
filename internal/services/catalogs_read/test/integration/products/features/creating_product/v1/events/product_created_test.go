//go:build integration
// +build integration

package events

import (
	"context"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/go-playground/validator"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/suite"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/messaging/types"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/test/messaging/consumer"
	testUtils "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/test/utils"
	externalEvents "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/read_service/internal/products/features/creating_product/v1/events/integration_events/external_events"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/read_service/internal/products/models"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/read_service/internal/shared/test_fixture/integration"
)

type productCreatedIntegrationTests struct {
	*integration.IntegrationTestSharedFixture
}

func TestCreateProductIntegration(t *testing.T) {
	suite.Run(
		t,
		&productCreatedIntegrationTests{
			IntegrationTestSharedFixture: integration.NewIntegrationTestSharedFixture(t),
		},
	)
}

func (c *productCreatedIntegrationTests) Test_Product_Created_Consumer_Should_Consume_Product_Created() {
	ctx := context.Background()
	fakeConsumer := consumer.NewRabbitMQFakeTestConsumerHandler[*externalEvents.ProductCreatedV1]()
	err := c.Bus.ConnectConsumerHandler(&externalEvents.ProductCreatedV1{}, fakeConsumer)
	c.NoError(err)

	// in test mode we set rabbitmq `AutoStart=false`, so we should run rabbitmq bus manually
	c.Bus.Start(context.Background())
	// wait for consumers ready to consume before publishing messages, preparation background workers takes a bit time (for preventing messages lost)
	time.Sleep(1 * time.Second)
	defer c.Bus.Stop()

	err = c.Bus.PublishMessage(
		ctx,
		&externalEvents.ProductCreatedV1{
			Message:   types.NewMessage(uuid.NewV4().String()),
			ProductId: uuid.NewV4().String(),
			Name:      gofakeit.Name(),
			Price:     gofakeit.Price(150, 6000),
			CreatedAt: time.Now(),
		},
		nil,
	)
	c.NoError(err)

	// ensuring message published to the rabbitmq broker
	c.NoError(testUtils.WaitUntilConditionMet(func() bool {
		return fakeConsumer.IsHandled()
	}))
}

func (c *productCreatedIntegrationTests) Test_Product_Created_Consumer() {
	ctx := context.Background()
	cons := externalEvents.NewProductCreatedConsumer(c.Log, validator.New(), c.Tracer)
	err := c.Bus.ConnectConsumerHandler(externalEvents.ProductCreatedV1{}, cons)
	c.NoError(err)

	// in test mode we set rabbitmq `AutoStart=false`, so we should run rabbitmq bus manually
	c.Bus.Start(context.Background())
	// wait for consumers ready to consume before publishing messages, preparation background workers takes a bit time (for preventing messages lost)
	time.Sleep(1 * time.Second)
	defer c.Bus.Stop()

	pid := uuid.NewV4().String()
	productCreated := &externalEvents.ProductCreatedV1{
		Message:     types.NewMessage(uuid.NewV4().String()),
		ProductId:   pid,
		CreatedAt:   time.Now(),
		Name:        gofakeit.Name(),
		Price:       gofakeit.Price(150, 6000),
		Description: gofakeit.AdjectiveDescriptive(),
	}

	err = c.Bus.PublishMessage(ctx, productCreated, nil)
	c.NoError(err)

	var p *models.Product

	c.NoError(testUtils.WaitUntilConditionMet(func() bool {
		p, err = c.ProductRepository.GetProductByProductId(ctx, pid)
		c.NoError(err)

		return p != nil
	}))

	c.NoError(err)
	c.NotNil(p)
	c.Equal(pid, p.ProductId)
}
