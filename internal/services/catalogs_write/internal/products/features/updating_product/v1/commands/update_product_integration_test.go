//go:build.sh integration
// +build.sh integration

package commands

import (
	"net/http"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/mehdihadeli/go-mediatr"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	customErrors "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/http_errors/custom_errors"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/test/messaging"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/test/messaging/consumer"
	testUtils "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/test/utils"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/write_service/internal/products/features/updating_product/v1/events/integration_events"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/write_service/internal/products/mocks/testData"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/write_service/internal/shared/test_fixtures/integration"
)

type updateProductIntegrationTests struct {
	*integration.IntegrationTestFixture
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
	testUtils.SkipCI(c.T())

	existing := testData.Products[0]

	command, err := NewUpdateProduct(
		existing.ProductId,
		gofakeit.Name(),
		existing.Description,
		existing.Price,
	)
	c.Require().NoError(err)

	result, err := mediatr.Send[*UpdateProduct, *mediatr.Unit](c.Ctx, command)
	c.Require().NoError(err)

	c.NotNil(result)

	updatedProduct, err := c.IntegrationTestFixture.ProductRepository.GetProductById(
		c.Ctx,
		existing.ProductId,
	)
	c.NotNil(updatedProduct)
	c.Equal(existing.ProductId, updatedProduct.ProductId)
	c.Equal(existing.Price, updatedProduct.Price)
	c.NotEqual(existing.Name, updatedProduct.Name)
}

func (c *updateProductIntegrationTests) Test_Should_Return_NotFound_Error_When_Item_DoesNot_Exist() {
	testUtils.SkipCI(c.T())

	id := uuid.NewV4()

	command, err := NewUpdateProduct(
		id,
		gofakeit.Name(),
		gofakeit.EmojiDescription(),
		gofakeit.Price(150, 6000),
	)
	c.Require().NoError(err)

	result, err := mediatr.Send[*UpdateProduct, *mediatr.Unit](c.Ctx, command)

	c.Assert().Error(err)
	c.True(customErrors.IsApplicationError(err, http.StatusNotFound))
	c.Assert().Nil(result)
}

func (c *updateProductIntegrationTests) Test_Should_Publish_Product_Updated_To_Broker() {
	testUtils.SkipCI(c.T())

	shouldPublish := messaging.ShouldProduced[*integration_events.ProductUpdatedV1](
		c.Ctx,
		c.Bus,
		nil,
	)

	existing := testData.Products[0]

	command, err := NewUpdateProduct(
		existing.ProductId,
		gofakeit.Name(),
		existing.Description,
		existing.Price,
	)
	c.Require().NoError(err)

	_, err = mediatr.Send[*UpdateProduct, *mediatr.Unit](c.Ctx, command)
	c.Require().NoError(err)

	// ensuring message published to the rabbitmq broker
	shouldPublish.Validate(c.Ctx, "there is no published message", time.Second*30)
}

func (c *updateProductIntegrationTests) Test_Should_Consume_Product_Created_With_Existing_Consumer_From_Broker() {
	testUtils.SkipCI(c.T())

	// should consume productCreatedTestConsumer
	newConsumer := messaging.ShouldConsume[*integration_events.ProductUpdatedV1](c.Ctx, c.Bus, nil)

	existing := testData.Products[0]
	command, err := NewUpdateProduct(
		existing.ProductId,
		gofakeit.Name(),
		existing.Description,
		existing.Price,
	)
	c.Require().NoError(err)

	_, err = mediatr.Send[*UpdateProduct, *mediatr.Unit](c.Ctx, command)
	c.Require().NoError(err)

	// ensuring message can be consumed with a consumer
	newConsumer.Validate(c.Ctx, "there is no consumed message", time.Second*30)
}

func (c *updateProductIntegrationTests) Test_Should_Consume_Product_Updated_With_New_Consumer_From_Broker() {
	testUtils.SkipCI(c.T())

	// should consume productCreatedTestConsumer
	newConsumer, err := messaging.ShouldConsumeNewConsumer[*integration_events.ProductUpdatedV1](
		c.Ctx,
		c.Bus,
	)
	require.NoError(c.T(), err)

	c.IntegrationTestFixture.Run()

	existing := testData.Products[0]
	command, err := NewUpdateProduct(
		existing.ProductId,
		gofakeit.Name(),
		existing.Description,
		existing.Price,
	)
	c.Require().NoError(err)

	_, err = mediatr.Send[*UpdateProduct, *mediatr.Unit](c.Ctx, command)
	c.Require().NoError(err)

	// ensuring message can be consumed with a consumer
	newConsumer.Validate(c.Ctx, "there is no consumed message", time.Second*30)
}

func (c *updateProductIntegrationTests) BeforeTest(suiteName, testName string) {
	if testName != "Test_Should_Consume_Product_Updated_With_New_Consumer_From_Broker" {
		c.IntegrationTestFixture.Run()
	}
}

func (c *updateProductIntegrationTests) SetupTest() {
	c.T().Log("SetupTest")
	c.IntegrationTestFixture = integration.NewIntegrationTestFixture(c.IntegrationTestSharedFixture)
	err := mediatr.RegisterRequestHandler[*UpdateProduct, *mediatr.Unit](
		NewUpdateProductHandler(c.Log, c.Cfg, c.CatalogUnitOfWorks, c.Bus),
	)
	c.Require().NoError(err)

	testConsumer := consumer.NewRabbitMQFakeTestConsumerHandler[*integration_events.ProductUpdatedV1]()
	err = c.Bus.ConnectConsumerHandler(&integration_events.ProductUpdatedV1{}, testConsumer)
	c.Require().NoError(err)
}

func (c *updateProductIntegrationTests) TearDownTest() {
	c.T().Log("TearDownTest")
	// cleanup test containers with their hooks
}
