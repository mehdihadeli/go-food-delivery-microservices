package commands

import (
	"net/http"
	"testing"
	"time"

	"github.com/mehdihadeli/go-mediatr"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/suite"

    customErrors "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/http_errors/custom_errors"
    "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/test/messaging"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/test/messaging/consumer"
    testUtils "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/test/utils"

    integrationEvents "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/write_service/internal/products/features/deleting_product/v1/events/integration_events"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/write_service/internal/products/mocks/testData"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/write_service/internal/shared/test_fixtures/integration"
)

type deleteProductIntegrationTests struct {
	*integration.IntegrationTestFixture
	*integration.IntegrationTestSharedFixture
}

func TestDeleteProductIntegration(t *testing.T) {
	suite.Run(t, &deleteProductIntegrationTests{IntegrationTestSharedFixture: integration.NewIntegrationTestSharedFixture(t)})
}

func (c *deleteProductIntegrationTests) Test_Should_Delete_Product_From_DB() {
	testUtils.SkipCI(c.T())

	id := testData.Products[0].ProductId
	command, err := NewDeleteProduct(id)
	c.Require().NoError(err)

	result, err := mediatr.Send[*DeleteProduct, *mediatr.Unit](c.Ctx, command)

	c.Require().NoError(err)
	c.Assert().NotNil(result)

	deletedProduct, err := c.IntegrationTestFixture.ProductRepository.GetProductById(c.Ctx, id)
	c.Assert().Nil(deletedProduct)
}

func (c *deleteProductIntegrationTests) Test_Should_Returns_NotFound_Error_When_Record_DoesNot_Exists() {
	testUtils.SkipCI(c.T())

	id := uuid.NewV4()
	command, err := NewDeleteProduct(id)
	c.Require().NoError(err)

	result, err := mediatr.Send[*DeleteProduct, *mediatr.Unit](c.Ctx, command)

	c.Assert().Error(err)
	c.True(customErrors.IsApplicationError(err, http.StatusNotFound))
	c.True(customErrors.IsNotFoundError(err))
	c.Assert().Nil(result)
}

func (c *deleteProductIntegrationTests) Test_Should_Publish_Product_Created_To_Broker() {
	testUtils.SkipCI(c.T())

	shouldPublish := messaging.ShouldProduced[*integrationEvents.ProductDeletedV1](c.Ctx, c.Bus, nil)

	id := testData.Products[0].ProductId
	command, err := NewDeleteProduct(id)
	c.Require().NoError(err)

	_, err = mediatr.Send[*DeleteProduct, *mediatr.Unit](c.Ctx, command)
	c.Require().NoError(err)

	// ensuring message published to the rabbitmq broker
	shouldPublish.Validate(c.Ctx, "there is no published message", time.Second*30)
}

func (c *deleteProductIntegrationTests) SetupTest() {
	c.T().Log("SetupTest")
	c.IntegrationTestFixture = integration.NewIntegrationTestFixture(c.IntegrationTestSharedFixture)
	err := mediatr.RegisterRequestHandler[*DeleteProduct, *mediatr.Unit](NewDeleteProductHandler(c.Log, c.Cfg, c.CatalogUnitOfWorks, c.Bus))
	c.Require().NoError(err)

	testConsumer := consumer.NewRabbitMQFakeTestConsumerHandlerWithHypothesis[*integrationEvents.ProductDeletedV1](nil)
	err = c.Bus.ConnectConsumerHandler(&integrationEvents.ProductDeletedV1{}, testConsumer)
	c.Require().NoError(err)

	c.IntegrationTestFixture.Run()
}

func (c *deleteProductIntegrationTests) TearDownTest() {
	c.T().Log("TearDownTest")
	// cleanup test containers with their hooks
}
