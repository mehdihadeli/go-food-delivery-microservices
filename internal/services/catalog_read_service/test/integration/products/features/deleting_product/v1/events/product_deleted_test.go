//go:build integration
// +build integration

package events

import (
	"context"
	"testing"
	"time"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/messaging/types"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/test/messaging"
	testUtils "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/test/utils"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/suite"

	externalEvents "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogreadservice/internal/products/features/deleting_products/v1/events/integration_events/external_events"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogreadservice/internal/products/models"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogreadservice/internal/shared/test_fixture/integration"
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

	// check for consuming `ProductDeletedV1` message with existing consumer
	hypothesis := messaging.ShouldConsume[*externalEvents.ProductDeletedV1](ctx, c.Bus, nil)

	// in test mode we set rabbitmq `AutoStart=false`, so we should run rabbitmq bus manually
	c.Bus.Start(context.Background())
	// wait for consumers ready to consume before publishing messages, preparation background workers takes a bit time (for preventing messages lost)
	time.Sleep(1 * time.Second)
	defer c.Bus.Stop()

	event := &externalEvents.ProductDeletedV1{
		Message:   types.NewMessage(uuid.NewV4().String()),
		ProductId: c.Items[0].ProductId,
	}

	err := c.Bus.PublishMessage(
		ctx,
		event,
		nil,
	)
	c.NoError(err)

	// ensuring message can be consumed with a consumer
	hypothesis.Validate(ctx, "there is no consumed message", time.Second*60)
}

func (c *productDeletedIntegrationTests) Test_Product_Deleted_Consumer_Should_Consume_Product_Deleted_With_New_Consumer() {
	ctx := context.Background()

	// check for consuming `ProductDeletedV1` message, with a new consumer
	hypothesis, err := messaging.ShouldConsumeNewConsumer[*externalEvents.ProductDeletedV1](c.Bus)
	c.Require().NoError(err)

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

	// ensuring message can be consumed with a consumer
	hypothesis.Validate(ctx, "there is no consumed message", time.Second*60)
}

func (c *productDeletedIntegrationTests) Test_Product_Deleted_Consumer() {
	ctx := context.Background()

	// in test mode we set rabbitmq `AutoStart=false`, so we should run rabbitmq bus manually
	c.Bus.Start(ctx)
	// wait for consumers ready to consume before publishing messages, preparation background workers takes a bit time (for preventing messages lost)
	time.Sleep(3 * time.Second)
	defer c.Bus.Stop()

	productDeleted := &externalEvents.ProductDeletedV1{
		Message:   types.NewMessage(uuid.NewV4().String()),
		ProductId: c.Items[0].ProductId,
	}

	err := c.Bus.PublishMessage(ctx, productDeleted, nil)
	c.NoError(err)

	var p *models.Product

	c.NoError(testUtils.WaitUntilConditionMet(func() bool {
		p, err = c.ProductRepository.GetProductByProductId(ctx, c.Items[0].ProductId)
		c.NoError(err)

		return p == nil
	}))

	c.Nil(p)
}
