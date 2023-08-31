//go:build integration
// +build integration

package v1

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	customErrors "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/http_errors/custom_errors"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/test/messaging"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/test/messaging/consumer"
	"github.com/mehdihadeli/go-mediatr"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/features/updating_product/v1/commands"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/features/updating_product/v1/events/integration_events"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/shared/test_fixtures/integration"
)

type updateProductIntegrationTests struct {
	*integration.IntegrationTestSharedFixture
}

func TestUpdateProductIntegration(t *testing.T) {
	suite.Run(
		t,
		&updateProductIntegrationTests{
			IntegrationTestSharedFixture: integration.NewIntegrationTestSharedFixture(t),
		},
	)
}

func (c *updateProductIntegrationTests) Test_Should_Update_Existing_Product_In_DB() {
	ctx := context.Background()
	existing := c.Items[0]

	command, err := commands.NewUpdateProduct(
		existing.ProductId,
		gofakeit.Name(),
		existing.Description,
		existing.Price,
	)
	c.Require().NoError(err)

	result, err := mediatr.Send[*commands.UpdateProduct, *mediatr.Unit](ctx, command)
	c.Require().NoError(err)

	c.NotNil(result)

	updatedProduct, err := c.ProductRepository.GetProductById(
		ctx,
		existing.ProductId,
	)
	c.NotNil(updatedProduct)
	c.Equal(existing.ProductId, updatedProduct.ProductId)
	c.Equal(existing.Price, updatedProduct.Price)
	c.NotEqual(existing.Name, updatedProduct.Name)
}

func (c *updateProductIntegrationTests) Test_Should_Return_NotFound_Error_When_Item_DoesNot_Exist() {
	ctx := context.Background()

	id := uuid.NewV4()

	command, err := commands.NewUpdateProduct(
		id,
		gofakeit.Name(),
		gofakeit.EmojiDescription(),
		gofakeit.Price(150, 6000),
	)
	c.Require().NoError(err)

	result, err := mediatr.Send[*commands.UpdateProduct, *mediatr.Unit](ctx, command)

	c.Assert().Error(err)
	c.True(customErrors.IsApplicationError(err, http.StatusNotFound))
	c.Assert().Nil(result)
}

func (c *updateProductIntegrationTests) Test_Should_Publish_Product_Updated_To_Broker() {
	ctx := context.Background()

	shouldPublish := messaging.ShouldProduced[*integration_events.ProductUpdatedV1](
		ctx,
		c.Bus,
		nil,
	)

	existing := c.Items[0]

	command, err := commands.NewUpdateProduct(
		existing.ProductId,
		gofakeit.Name(),
		existing.Description,
		existing.Price,
	)
	c.Require().NoError(err)

	_, err = mediatr.Send[*commands.UpdateProduct, *mediatr.Unit](ctx, command)
	c.Require().NoError(err)

	// ensuring message published to the rabbitmq broker
	shouldPublish.Validate(ctx, "there is no published message", time.Second*30)
}

func (c *updateProductIntegrationTests) Test_Should_Consume_Product_Created_With_Existing_Consumer_From_Broker() {
	ctx := context.Background()

	// we don't have a consumer in this service, so we simulate one consumer in `SetupSuite`
	// // check for consuming `ProductUpdatedV1` message with existing consumer
	hypothesis := messaging.ShouldConsume[*integration_events.ProductUpdatedV1](ctx, c.Bus, nil)

	existing := c.Items[0]
	command, err := commands.NewUpdateProduct(
		existing.ProductId,
		gofakeit.Name(),
		existing.Description,
		existing.Price,
	)
	c.Require().NoError(err)

	_, err = mediatr.Send[*commands.UpdateProduct, *mediatr.Unit](ctx, command)
	c.Require().NoError(err)

	// ensuring message can be consumed with a consumer
	hypothesis.Validate(ctx, "there is no consumed message", time.Second*30)
}

func (c *updateProductIntegrationTests) Test_Should_Consume_Product_Updated_With_New_Consumer_From_Broker() {
	ctx := context.Background()

	//  check for consuming `ProductUpdatedV1` message, with a new consumer
	hypothesis, err := messaging.ShouldConsumeNewConsumer[*integration_events.ProductUpdatedV1](
		c.Bus,
	)
	require.NoError(c.T(), err)

	// at first, we should add new consumer to rabbitmq bus then start the broker, because we can't add new consumer after start.
	// we should also turn off consumer in `BeforeTest` for this test
	c.Bus.Start(ctx)

	// wait for consumers ready to consume before publishing messages, preparation background workers takes a bit time (for preventing messages lost)
	time.Sleep(1 * time.Second)

	existing := c.Items[0]
	command, err := commands.NewUpdateProduct(
		existing.ProductId,
		gofakeit.Name(),
		existing.Description,
		existing.Price,
	)
	c.Require().NoError(err)

	_, err = mediatr.Send[*commands.UpdateProduct, *mediatr.Unit](ctx, command)
	c.Require().NoError(err)

	// ensuring message can be consumed with a consumer
	hypothesis.Validate(ctx, "there is no consumed message", time.Second*30)
}

func (c *updateProductIntegrationTests) BeforeTest(suiteName, testName string) {
	if testName == "Test_Should_Consume_Product_Updated_With_New_Consumer_From_Broker" {
		c.Bus.Stop()
	}
}

func (c *updateProductIntegrationTests) SetupSuite() {
	// we don't have a consumer in this service, so we simulate one consumer, register one consumer for `ProductUpdatedV1` message before executing the tests
	productUpdatedConsumer := consumer.NewRabbitMQFakeTestConsumerHandler[*integration_events.ProductUpdatedV1]()
	err := c.Bus.ConnectConsumerHandler(
		&integration_events.ProductUpdatedV1{},
		productUpdatedConsumer,
	)
	c.Require().NoError(err)

	// in test mode we set rabbitmq `AutoStart=false`, so we should run rabbitmq bus manually
	c.Bus.Start(context.Background())
	// wait for consumers ready to consume before publishing messages, preparation background workers takes a bit time (for preventing messages lost)
	time.Sleep(1 * time.Second)
}

func (c *updateProductIntegrationTests) TearDownSuite() {
	c.Bus.Stop()
}
