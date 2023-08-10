//go:build integration
// +build integration

package events

import (
	"context"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/messaging/types"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/test/messaging"
	testUtils "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/test/utils"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/suite"

	externalEvents "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogreadservice/internal/products/features/updating_products/v1/events/integration_events/external_events"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogreadservice/internal/products/models"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogreadservice/internal/shared/test_fixture/integration"
)

type productUpdatedIntegrationTests struct {
	*integration.IntegrationTestSharedFixture
}

func TestProductUpdatedIntegration(t *testing.T) {
	suite.Run(
		t,
		&productUpdatedIntegrationTests{
			IntegrationTestSharedFixture: integration.NewIntegrationTestSharedFixture(t),
		},
	)
}

func (c *productUpdatedIntegrationTests) Test_Product_Updated_Consumer_Should_Consume_Product_Updated() {
	ctx := context.Background()
	// check for consuming `ProductUpdatedV1` message with existing consumer
	hypothesis := messaging.ShouldConsume[*externalEvents.ProductUpdatedV1](ctx, c.Bus, nil)

	// in test mode we set rabbitmq `AutoStart=false`, so we should run rabbitmq bus manually
	c.Bus.Start(context.Background())
	time.Sleep(1 * time.Second)
	defer c.Bus.Stop()
	// wait for consumers ready to consume before publishing messages, preparation background workers takes a bit time (for preventing messages lost)

	err := c.Bus.PublishMessage(
		ctx,
		&externalEvents.ProductUpdatedV1{
			Message:     types.NewMessage(uuid.NewV4().String()),
			ProductId:   c.Items[0].ProductId,
			Name:        gofakeit.Name(),
			Price:       gofakeit.Price(100, 1000),
			Description: gofakeit.EmojiDescription(),
			UpdatedAt:   time.Now(),
		},
		nil,
	)
	c.NoError(err)

	// ensuring message can be consumed with a consumer
	hypothesis.Validate(ctx, "there is no consumed message", time.Second*30)
}

func (c *productUpdatedIntegrationTests) Test_Product_Updated_Consumer_Should_Consume_Product_Created_With_New_Consumer() {
	ctx := context.Background()

	// check for consuming `ProductUpdatedV1` message, with a new consumer
	hypothesis, err := messaging.ShouldConsumeNewConsumer[*externalEvents.ProductUpdatedV1](c.Bus)

	// in test mode we set rabbitmq `AutoStart=false`, so we should run rabbitmq bus manually
	c.Bus.Start(context.Background())
	time.Sleep(1 * time.Second)
	defer c.Bus.Stop()
	// wait for consumers ready to consume before publishing messages, preparation background workers takes a bit time (for preventing messages lost)

	err = c.Bus.PublishMessage(
		ctx,
		&externalEvents.ProductUpdatedV1{
			Message:     types.NewMessage(uuid.NewV4().String()),
			ProductId:   c.Items[0].ProductId,
			Name:        gofakeit.Name(),
			Price:       gofakeit.Price(100, 1000),
			Description: gofakeit.EmojiDescription(),
			UpdatedAt:   time.Now(),
		},
		nil,
	)
	c.NoError(err)

	// ensuring message can be consumed with a consumer
	hypothesis.Validate(ctx, "there is no consumed message", time.Second*30)
}

func (c *productUpdatedIntegrationTests) Test_Product_Updated_Consumer() {
	ctx := context.Background()

	// in test mode we set rabbitmq `AutoStart=false`, so we should run rabbitmq bus manually
	c.Bus.Start(context.Background())
	// wait for consumers ready to consume before publishing messages, preparation background workers takes a bit time (for preventing messages lost)
	time.Sleep(1 * time.Second)
	defer c.Bus.Stop()

	productUpdated := &externalEvents.ProductUpdatedV1{
		Message:     types.NewMessage(uuid.NewV4().String()),
		ProductId:   c.Items[0].ProductId,
		Name:        gofakeit.Name(),
		Description: gofakeit.AdjectiveDescriptive(),
		Price:       gofakeit.Price(150, 6000),
		UpdatedAt:   time.Now(),
	}

	err := c.Bus.PublishMessage(ctx, productUpdated, nil)
	c.NoError(err)

	var p *models.Product

	c.NoError(testUtils.WaitUntilConditionMet(func() bool {
		p, err = c.ProductRepository.GetProductByProductId(ctx, c.Items[0].ProductId)

		return p != nil && p.Name == productUpdated.Name
	}))

	c.NotNil(p)
	c.Equal(productUpdated.Name, p.Name)
	c.Equal(productUpdated.ProductId, p.ProductId)
}
