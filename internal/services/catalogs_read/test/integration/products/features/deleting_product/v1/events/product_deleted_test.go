//go:build integration
// +build integration

package events

import (
	"context"
	"testing"
	"time"

	"github.com/go-playground/validator"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/suite"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/messaging/types"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/test/messaging/consumer"
	testUtils "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/test/utils"
	externalEvents "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/read_service/internal/products/features/deleting_products/v1/events/integration_events/external_events"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/read_service/internal/shared/test_fixture/integration"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/read_service/internal/products/models"
)

type productDeletedIntegrationTests struct {
	*integration.IntegrationTestSharedFixture
}

func TestProductDeletedIntegration(t *testing.T) {
	suite.Run(
		t,
		&productDeletedIntegrationTests{
			IntegrationTestSharedFixture: integration.NewIntegrationTestSharedFixture(t),
		},
	)
}

func (c *productDeletedIntegrationTests) Test_Product_Deleted_Consumer_Should_Consume_Product_Deleted() {
	ctx := context.Background()

	fakeConsumer := consumer.NewRabbitMQFakeTestConsumerHandler[externalEvents.ProductDeletedV1]()
	err := c.Bus.ConnectConsumerHandler(externalEvents.ProductDeletedV1{}, fakeConsumer)
	c.NoError(err)

	// in test mode we set rabbitmq `AutoStart=false`, so we should run rabbitmq bus manually
	c.Bus.Start(context.Background())
	// wait for consumers ready to consume before publishing messages, preparation background workers takes a bit time (for preventing messages lost)
	time.Sleep(1 * time.Second)
	defer c.Bus.Stop()

	event := &externalEvents.ProductDeletedV1{
		Message:   types.NewMessage(uuid.NewV4().String()),
		ProductId: c.Items[0].ProductId,
	}

	err = c.Bus.PublishMessage(
		ctx,
		event,
		nil,
	)
	c.NoError(err)

	// ensuring message published to the rabbitmq broker
	c.NoError(testUtils.WaitUntilConditionMet(func() bool {
		return fakeConsumer.IsHandled()
	}))
}

func (c *productDeletedIntegrationTests) Test_Product_Deleted_Consumer() {
	ctx := context.Background()

	cons := externalEvents.NewProductDeletedConsumer(c.Log, validator.New(), c.Tracer)
	err := c.Bus.ConnectConsumerHandler(externalEvents.ProductDeletedV1{}, cons)
	c.NoError(err)

	// in test mode we set rabbitmq `AutoStart=false`, so we should run rabbitmq bus manually
	c.Bus.Start(ctx)
	// wait for consumers ready to consume before publishing messages, preparation background workers takes a bit time (for preventing messages lost)
	time.Sleep(3 * time.Second)
	defer c.Bus.Stop()

	pid := "ff13c422-e0dc-466d-9bee-d09c1d3122e1"
	productDeleted := &externalEvents.ProductDeletedV1{
		Message:   types.NewMessage(uuid.NewV4().String()),
		ProductId: pid,
	}

	err = c.Bus.PublishMessage(ctx, productDeleted, nil)
	c.NoError(err)

	var p *models.Product

	c.NoError(testUtils.WaitUntilConditionMet(func() bool {
		p, err = c.ProductRepository.GetProductByProductId(ctx, pid)
		c.NoError(err)

		return p == nil
	}))

	c.Nil(p)
}
