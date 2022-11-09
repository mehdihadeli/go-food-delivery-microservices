package v1

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/test/messaging"

	customErrors "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/http_errors/custom_errors"

	uuid "github.com/satori/go.uuid"

	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/updating_product/events/integration/v1"

	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/test/utils"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/mocks/testData"
	"github.com/stretchr/testify/suite"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/mehdihadeli/go-mediatr"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/test/messaging/consumer"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/shared/test_fixtures/integration"
)

type updateProductIntegrationTests struct {
	*integration.IntegrationTestFixture
	*integration.IntegrationTestSharedFixture
}

func TestUpdateProductIntegration(t *testing.T) {
	suite.Run(t, &updateProductIntegrationTests{IntegrationTestSharedFixture: integration.NewIntegrationTestSharedFixture(t)})
}

func (c *updateProductIntegrationTests) Test_Should_Update_Existing_Product_In_DB() {
	utils.SkipCI(c.T())

	existing := testData.Products[0]

	command := NewUpdateProduct(existing.ProductId, gofakeit.Name(), existing.Description, existing.Price)

	result, err := mediatr.Send[*UpdateProduct, *mediatr.Unit](c.Ctx, command)
	c.Require().NoError(err)

	c.NotNil(result)

	updatedProduct, err := c.IntegrationTestFixture.ProductRepository.GetProductById(c.Ctx, existing.ProductId)
	c.NotNil(updatedProduct)
	c.Equal(existing.ProductId, updatedProduct.ProductId)
	c.Equal(existing.Price, updatedProduct.Price)
	c.NotEqual(existing.Name, updatedProduct.Name)
}

func (c *updateProductIntegrationTests) Test_Should_Return_NotFound_Error_When_Item_DoesNot_Exist() {
	utils.SkipCI(c.T())

	id := uuid.NewV4()

	command := NewUpdateProduct(id, gofakeit.Name(), gofakeit.EmojiDescription(), gofakeit.Price(150, 6000))

	result, err := mediatr.Send[*UpdateProduct, *mediatr.Unit](c.Ctx, command)

	c.Assert().Error(err)
	c.True(customErrors.IsApplicationError(err, http.StatusNotFound))
	c.Assert().Nil(result)
}

func (c *updateProductIntegrationTests) Test_Should_Publish_Product_Updated_To_Broker() {
	utils.SkipCI(c.T())

	shouldPublish := messaging.ShouldProduced[*v1.ProductUpdatedV1](c.Ctx, c.Bus, nil)

	existing := testData.Products[0]

	command := NewUpdateProduct(existing.ProductId, gofakeit.Name(), existing.Description, existing.Price)
	_, err := mediatr.Send[*UpdateProduct, *mediatr.Unit](c.Ctx, command)
	c.Require().NoError(err)

	// ensuring message published to the rabbitmq broker
	shouldPublish.Validate(c.Ctx, "there is no published message", time.Second*30)
}

func (c *updateProductIntegrationTests) Test_Should_Consume_Product_Created_With_Existing_Consumer_From_Broker() {
	utils.SkipCI(c.T())

	// should consume productCreatedTestConsumer
	newConsumer := messaging.ShouldConsume[*v1.ProductUpdatedV1](c.Ctx, c.Bus, nil)

	existing := testData.Products[0]
	command := NewUpdateProduct(existing.ProductId, gofakeit.Name(), existing.Description, existing.Price)
	_, err := mediatr.Send[*UpdateProduct, *mediatr.Unit](c.Ctx, command)
	c.Require().NoError(err)

	// ensuring message can be consumed with a consumer
	newConsumer.Validate(c.Ctx, "there is no consumed message", time.Second*30)
}

func (c *updateProductIntegrationTests) Test_Should_Consume_Product_Updated_With_New_Consumer_From_Broker() {
	utils.SkipCI(c.T())

	// should consume productCreatedTestConsumer
	newConsumer, err := messaging.ShouldConsumeNewConsumer[*v1.ProductUpdatedV1](c.Ctx, c.Bus)
	require.NoError(c.T(), err)

	c.IntegrationTestFixture.Run()

	existing := testData.Products[0]
	command := NewUpdateProduct(existing.ProductId, gofakeit.Name(), existing.Description, existing.Price)
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
	err := mediatr.RegisterRequestHandler[*UpdateProduct, *mediatr.Unit](NewUpdateProductHandler(c.Log, c.Cfg, c.ProductRepository, c.Bus))
	c.Require().NoError(err)

	testConsumer := consumer.NewRabbitMQFakeTestConsumerHandler[*v1.ProductUpdatedV1]()
	err = c.Bus.ConnectConsumerHandler(&v1.ProductUpdatedV1{}, testConsumer)
	c.Require().NoError(err)
}

func (c *updateProductIntegrationTests) TearDownTest() {
	c.T().Log("TearDownTest")
	// cleanup test containers with their hooks
}
