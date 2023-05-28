//go:build.sh integration
// +build.sh integration

package createProductCommand

import (
	"net/http"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/mehdihadeli/go-mediatr"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	customErrors "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/http_errors/custom_errors"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/test/messaging"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/test/messaging/consumer"
	testUtils "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/test/utils"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/write_service/internal/products/features/creating_product/v1/dtos"
	integrationEvents "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/write_service/internal/products/features/creating_product/v1/events/integration_events"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/write_service/internal/products/mocks/testData"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/write_service/internal/shared/test_fixtures/integration"
)

type createProductIntegrationTests struct {
	*integration.IntegrationTestFixture
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
	testUtils.SkipCI(c.T())

	command, err := NewCreateProduct(
		gofakeit.Name(),
		gofakeit.AdjectiveDescriptive(),
		gofakeit.Price(150, 6000),
	)
	c.Require().NoError(err)

	result, err := mediatr.Send[*CreateProduct, *dtos.CreateProductResponseDto](c.Ctx, command)
	c.Require().NoError(err)

	c.Assert().NotNil(result)
	c.Assert().Equal(command.ProductID, result.ProductID)

	createdProduct, err := c.IntegrationTestFixture.ProductRepository.GetProductById(
		c.Ctx,
		result.ProductID,
	)
	c.Require().NoError(err)
	c.Assert().NotNil(createdProduct)
}

func (c *createProductIntegrationTests) Test_Should_Return_Error_For_Duplicate_Record() {
	testUtils.SkipCI(c.T())

	id := testData.Products[0].ProductId

	command := &CreateProduct{
		Name:        gofakeit.Name(),
		Description: gofakeit.AdjectiveDescriptive(),
		Price:       gofakeit.Price(150, 6000),
		ProductID:   id,
	}

	result, err := mediatr.Send[*CreateProduct, *dtos.CreateProductResponseDto](c.Ctx, command)
	c.Assert().Error(err)
	c.True(customErrors.IsApplicationError(err, http.StatusConflict))
	c.Assert().Nil(result)
}

func (c *createProductIntegrationTests) Test_Should_Publish_Product_Created_To_Broker() {
	testUtils.SkipCI(c.T())

	shouldPublish := messaging.ShouldProduced[*integrationEvents.ProductCreatedV1](
		c.Ctx,
		c.Bus,
		nil,
	)

	command, err := NewCreateProduct(
		gofakeit.Name(),
		gofakeit.AdjectiveDescriptive(),
		gofakeit.Price(150, 6000),
	)
	c.Require().NoError(err)

	_, err = mediatr.Send[*CreateProduct, *dtos.CreateProductResponseDto](c.Ctx, command)
	c.Require().NoError(err)

	// ensuring message published to the rabbitmq broker
	shouldPublish.Validate(c.Ctx, "there is no published message", time.Second*30)
}

func (c *createProductIntegrationTests) Test_Should_Consume_Product_Created_With_Existing_Consumer_From_Broker() {
	testUtils.SkipCI(c.T())

	// should consume productCreatedTestConsumer
	newConsumer := messaging.ShouldConsume[*integrationEvents.ProductCreatedV1](c.Ctx, c.Bus, nil)

	command, err := NewCreateProduct(
		gofakeit.Name(),
		gofakeit.AdjectiveDescriptive(),
		gofakeit.Price(150, 6000),
	)
	c.Require().NoError(err)

	_, err = mediatr.Send[*CreateProduct, *dtos.CreateProductResponseDto](c.Ctx, command)
	c.Require().NoError(err)

	// ensuring message can be consumed with a consumer
	newConsumer.Validate(c.Ctx, "there is no consumed message", time.Second*30)
}

func (c *createProductIntegrationTests) Test_Should_Consume_Product_Created_With_New_Consumer_From_Broker() {
	testUtils.SkipCI(c.T())

	// should consume productCreatedTestConsumer
	newConsumer, err := messaging.ShouldConsumeNewConsumer[*integrationEvents.ProductCreatedV1](
		c.Ctx,
		c.Bus,
	)
	require.NoError(c.T(), err)

	c.IntegrationTestFixture.Run()

	command, err := NewCreateProduct(
		gofakeit.Name(),
		gofakeit.AdjectiveDescriptive(),
		gofakeit.Price(150, 6000),
	)
	c.Require().NoError(err)

	_, err = mediatr.Send[*CreateProduct, *dtos.CreateProductResponseDto](c.Ctx, command)
	c.Require().NoError(err)

	// ensuring message can be consumed with a consumer
	newConsumer.Validate(c.Ctx, "there is no consumed message", time.Second*30)
}

func (c *createProductIntegrationTests) BeforeTest(suiteName, testName string) {
	if testName != "Test_Should_Consume_Product_Created_With_New_Consumer_From_Broker" {
		c.IntegrationTestFixture.Run()
	}
}

func (c *createProductIntegrationTests) SetupTest() {
	c.T().Log("SetupTest")
	c.IntegrationTestFixture = integration.NewIntegrationTestFixture(c.IntegrationTestSharedFixture)
	err := mediatr.RegisterRequestHandler[*CreateProduct, *dtos.CreateProductResponseDto](
		NewCreateProductHandler(c.Log, c.Cfg, c.CatalogUnitOfWorks, c.Bus),
	)
	c.Require().NoError(err)

	testConsumer := consumer.NewRabbitMQFakeTestConsumerHandler[*integrationEvents.ProductCreatedV1]()
	err = c.Bus.ConnectConsumerHandler(&integrationEvents.ProductCreatedV1{}, testConsumer)
	c.Require().NoError(err)
}

func (c *createProductIntegrationTests) TearDownTest() {
	c.T().Log("TearDownTest")
	// cleanup test containers with their hooks
}
