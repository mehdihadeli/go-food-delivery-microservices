package v1

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/mehdihadeli/go-mediatr"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/test/messaging"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/test/messaging/consumer"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/test/utils"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/creating_product/dtos"
	v1 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/creating_product/events/integration/v1"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/shared/test_fixtures/integration"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type createProductIntegrationTest struct {
	*integration.IntegrationTestFixture
	*integration.IntegrationTestSharedFixture
}

func TestCreateProductIntegration(t *testing.T) {
	suite.Run(t, &createProductIntegrationTest{IntegrationTestSharedFixture: integration.NewIntegrationTestSharedFixture(t)})
}

func (c *createProductIntegrationTest) Test_Should_Create_New_Product_To_DB() {
	utils.SkipCI(c.T())

	command := NewCreateProduct(gofakeit.Name(), gofakeit.AdjectiveDescriptive(), gofakeit.Price(150, 6000))
	result, err := mediatr.Send[*CreateProduct, *dtos.CreateProductResponseDto](c.Ctx, command)
	c.Require().NoError(err)

	c.Assert().NotNil(result)
	c.Assert().Equal(command.ProductID, result.ProductID)

	createdProduct, err := c.IntegrationTestFixture.ProductRepository.GetProductById(c.Ctx, result.ProductID)
	c.Require().NoError(err)
	c.Assert().NotNil(createdProduct)
}

func (c *createProductIntegrationTest) Test_Should_Publish_Product_Created_To_Broker() {
	utils.SkipCI(c.T())

	shouldPublish := messaging.ShouldProduced[*v1.ProductCreatedV1](c.Ctx, c.Bus, nil)

	command := NewCreateProduct(gofakeit.Name(), gofakeit.AdjectiveDescriptive(), gofakeit.Price(150, 6000))
	_, err := mediatr.Send[*CreateProduct, *dtos.CreateProductResponseDto](c.Ctx, command)
	c.Require().NoError(err)

	// ensuring message published to the rabbitmq broker
	shouldPublish.Validate(c.Ctx, "there is no published message", time.Second*30)
}

func (c *createProductIntegrationTest) Test_Should_Consume_Product_Created_With_Existing_Consumer_From_Broker() {
	utils.SkipCI(c.T())

	// should consume productCreatedTestConsumer
	newConsumer := messaging.ShouldConsume[*v1.ProductCreatedV1](c.Ctx, c.Bus, nil)

	command := NewCreateProduct(gofakeit.Name(), gofakeit.AdjectiveDescriptive(), gofakeit.Price(150, 6000))
	_, err := mediatr.Send[*CreateProduct, *dtos.CreateProductResponseDto](c.Ctx, command)
	c.Require().NoError(err)

	// ensuring message can be consumed with a consumer
	newConsumer.Validate(c.Ctx, "there is no consumed message", time.Second*30)
}

func (c *createProductIntegrationTest) Test_Should_Consume_Product_Created_With_New_Consumer_From_Broker() {
	utils.SkipCI(c.T())

	// should consume productCreatedTestConsumer
	newConsumer, err := messaging.ShouldConsumeNewConsumer[*v1.ProductCreatedV1](c.Ctx, c.Bus)
	require.NoError(c.T(), err)

	c.IntegrationTestFixture.Run()

	command := NewCreateProduct(gofakeit.Name(), gofakeit.AdjectiveDescriptive(), gofakeit.Price(150, 6000))
	_, err = mediatr.Send[*CreateProduct, *dtos.CreateProductResponseDto](c.Ctx, command)
	c.Require().NoError(err)

	// ensuring message can be consumed with a consumer
	newConsumer.Validate(c.Ctx, "there is no consumed message", time.Second*30)
}

func (c *createProductIntegrationTest) BeforeTest(suiteName, testName string) {
	if testName != "Test_Should_Consume_Product_Created_With_New_Consumer_From_Broker" {
		c.IntegrationTestFixture.Run()
	}
}

func (c *createProductIntegrationTest) SetupTest() {
	c.T().Log("SetupTest")
	c.IntegrationTestFixture = integration.NewIntegrationTestFixture(c.IntegrationTestSharedFixture)
	err := mediatr.RegisterRequestHandler[*CreateProduct, *dtos.CreateProductResponseDto](NewCreateProductHandler(c.IntegrationTestSharedFixture.Log, c.IntegrationTestSharedFixture.Cfg, c.ProductRepository, c.Bus))
	c.Require().NoError(err)

	testConsumer := consumer.NewRabbitMQFakeTestConsumerHandler[*v1.ProductCreatedV1](nil)
	err = c.Bus.ConnectConsumerHandler(&v1.ProductCreatedV1{}, testConsumer)
	c.Require().NoError(err)
}

func (c *createProductIntegrationTest) TearDownTest() {
	c.T().Log("TearDownTest")
	// cleanup test containers with their hooks
}
