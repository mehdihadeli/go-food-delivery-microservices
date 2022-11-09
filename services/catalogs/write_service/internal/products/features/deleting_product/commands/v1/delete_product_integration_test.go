package v1

import (
	"net/http"
	"testing"
	"time"

	"github.com/mehdihadeli/go-mediatr"
	customErrors "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/http_errors/custom_errors"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/test/messaging"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/test/messaging/consumer"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/test/utils"
	v12 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/deleting_product/events/integration/v1"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/mocks/testData"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/shared/test_fixtures/integration"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/suite"
)

type deleteProductIntegrationTests struct {
	*integration.IntegrationTestFixture
	*integration.IntegrationTestSharedFixture
}

func TestDeleteProductIntegration(t *testing.T) {
	suite.Run(t, &deleteProductIntegrationTests{IntegrationTestSharedFixture: integration.NewIntegrationTestSharedFixture(t)})
}

func (c *deleteProductIntegrationTests) Test_Should_Delete_Product_From_DB() {
	utils.SkipCI(c.T())

	id := testData.Products[0].ProductId
	command := NewDeleteProduct(id)
	result, err := mediatr.Send[*DeleteProduct, *mediatr.Unit](c.Ctx, command)

	c.Require().NoError(err)
	c.Assert().NotNil(result)

	deletedProduct, err := c.IntegrationTestFixture.ProductRepository.GetProductById(c.Ctx, id)
	c.Assert().Nil(deletedProduct)
}

func (c *deleteProductIntegrationTests) Test_Should_Returns_NotFound_Error_When_Record_DoesNot_Exists() {
	utils.SkipCI(c.T())

	id := uuid.NewV4()
	command := NewDeleteProduct(id)
	result, err := mediatr.Send[*DeleteProduct, *mediatr.Unit](c.Ctx, command)

	c.Assert().Error(err)
	c.True(customErrors.IsApplicationError(err, http.StatusNotFound))
	c.True(customErrors.IsNotFoundError(err))
	c.Assert().Nil(result)
}

func (c *deleteProductIntegrationTests) Test_Should_Publish_Product_Created_To_Broker() {
	utils.SkipCI(c.T())

	shouldPublish := messaging.ShouldProduced[*v12.ProductDeletedV1](c.Ctx, c.Bus, nil)

	id := testData.Products[0].ProductId
	command := NewDeleteProduct(id)
	_, err := mediatr.Send[*DeleteProduct, *mediatr.Unit](c.Ctx, command)
	c.Require().NoError(err)

	// ensuring message published to the rabbitmq broker
	shouldPublish.Validate(c.Ctx, "there is no published message", time.Second*30)
}

func (c *deleteProductIntegrationTests) SetupTest() {
	c.T().Log("SetupTest")
	c.IntegrationTestFixture = integration.NewIntegrationTestFixture(c.IntegrationTestSharedFixture)
	err := mediatr.RegisterRequestHandler[*DeleteProduct, *mediatr.Unit](NewDeleteProductHandler(c.Log, c.Cfg, c.ProductRepository, c.Bus))
	c.Require().NoError(err)

	testConsumer := consumer.NewRabbitMQFakeTestConsumerHandlerWithHypothesis[*v12.ProductDeletedV1](nil)
	err = c.Bus.ConnectConsumerHandler(&v12.ProductDeletedV1{}, testConsumer)
	c.Require().NoError(err)

	c.IntegrationTestFixture.Run()
}

func (c *deleteProductIntegrationTests) TearDownTest() {
	c.T().Log("TearDownTest")
	// cleanup test containers with their hooks
}
