//go:build integration
// +build integration

package v1

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/mehdihadeli/go-mediatr"
	"github.com/stretchr/testify/suite"

	customErrors "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/http_errors/custom_errors"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/test/messaging"
	createProductCommand "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/write_service/internal/products/features/creating_product/v1/commands"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/write_service/internal/products/features/creating_product/v1/dtos"
	integrationEvents "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/write_service/internal/products/features/creating_product/v1/events/integration_events"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/write_service/internal/shared/test_fixtures/integration"
)

type createProductIntegrationTests struct {
	*integration.IntegrationTestSharedFixture
}

func TestCreateProductIntegration(t *testing.T) {
	suite.Run(
		t,
		&createProductIntegrationTests{
			IntegrationTestSharedFixture: integration.NewIntegrationTestSharedFixture(t),
		},
	)
}

func (c *createProductIntegrationTests) Test_Should_Create_New_Product_To_DB() {
	ctx := context.Background()

	command, err := createProductCommand.NewCreateProduct(
		gofakeit.Name(),
		gofakeit.AdjectiveDescriptive(),
		gofakeit.Price(150, 6000),
	)
	c.Require().NoError(err)

	result, err := mediatr.Send[*createProductCommand.CreateProduct, *dtos.CreateProductResponseDto](
		ctx,
		command,
	)
	c.Require().NoError(err)

	c.Assert().NotNil(result)
	c.Assert().Equal(command.ProductID, result.ProductID)

	createdProduct, err := c.ProductRepository.GetProductById(
		ctx,
		result.ProductID,
	)
	c.Require().NoError(err)
	c.Assert().NotNil(createdProduct)
}

func (c *createProductIntegrationTests) Test_Should_Return_Error_For_Duplicate_Record() {
	ctx := context.Background()

	id := c.Items[0].ProductId

	command := &createProductCommand.CreateProduct{
		Name:        gofakeit.Name(),
		Description: gofakeit.AdjectiveDescriptive(),
		Price:       gofakeit.Price(150, 6000),
		ProductID:   id,
	}

	result, err := mediatr.Send[*createProductCommand.CreateProduct, *dtos.CreateProductResponseDto](
		ctx,
		command,
	)
	c.Assert().Error(err)
	c.True(customErrors.IsApplicationError(err, http.StatusConflict))
	c.Assert().Nil(result)
}

func (c *createProductIntegrationTests) Test_Should_Publish_Product_Created_To_Broker() {
	ctx := context.Background()

	shouldPublish := messaging.ShouldProduced[*integrationEvents.ProductCreatedV1](
		ctx,
		c.Bus,
		nil,
	)

	command, err := createProductCommand.NewCreateProduct(
		gofakeit.Name(),
		gofakeit.AdjectiveDescriptive(),
		gofakeit.Price(150, 6000),
	)
	c.Require().NoError(err)

	_, err = mediatr.Send[*createProductCommand.CreateProduct, *dtos.CreateProductResponseDto](
		ctx,
		command,
	)
	c.Require().NoError(err)

	// ensuring message published to the rabbitmq broker
	shouldPublish.Validate(ctx, "there is no published message", time.Second*30)
}

func (c *createProductIntegrationTests) Test_Should_Consume_Product_Created_With_Existing_Consumer_From_Broker() {
	ctx := context.Background()

	// should consume productCreatedTestConsumer
	// we setup this handler in `BeforeTest`
	newConsumer := messaging.ShouldConsume[*integrationEvents.ProductCreatedV1](ctx, c.Bus, nil)

	command, err := createProductCommand.NewCreateProduct(
		gofakeit.Name(),
		gofakeit.AdjectiveDescriptive(),
		gofakeit.Price(150, 6000),
	)
	c.Require().NoError(err)

	_, err = mediatr.Send[*createProductCommand.CreateProduct, *dtos.CreateProductResponseDto](
		ctx,
		command,
	)
	c.Require().NoError(err)

	// ensuring message can be consumed with a consumer
	newConsumer.Validate(ctx, "there is no consumed message", time.Second*30)
}

func (c *createProductIntegrationTests) Test_Should_Consume_Product_Created_With_New_Consumer_From_Broker() {
	ctx := context.Background()
	defer c.Bus.Stop()

	// should consume productCreatedTestConsumer
	newConsumer, err := messaging.ShouldConsumeNewConsumer[*integrationEvents.ProductCreatedV1](
		c.Bus,
	)
	c.Require().NoError(err)

	// at first, we should add new consumer to rabbitmq bus then start the broker, because we can't add new consumer after start.
	// we should also turn off consumer in `BeforeTest` for this test
	c.Bus.Start(ctx)

	// wait for consumers ready to consume before publishing messages, preparation background workers takes a bit time (for preventing messages lost)
	time.Sleep(1 * time.Second)

	command, err := createProductCommand.NewCreateProduct(
		gofakeit.Name(),
		gofakeit.AdjectiveDescriptive(),
		gofakeit.Price(150, 6000),
	)
	c.Require().NoError(err)

	_, err = mediatr.Send[*createProductCommand.CreateProduct, *dtos.CreateProductResponseDto](
		ctx,
		command,
	)
	c.Require().NoError(err)

	// ensuring message can be consumed with a consumer
	newConsumer.Validate(ctx, "there is no consumed message", time.Second*30)
}

func (c *createProductIntegrationTests) BeforeTest(suiteName, testName string) {
	if testName == "Test_Should_Consume_Product_Created_With_New_Consumer_From_Broker" {
		c.Bus.Stop()
	}
}
