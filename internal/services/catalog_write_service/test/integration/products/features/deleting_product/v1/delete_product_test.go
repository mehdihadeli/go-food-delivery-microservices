//go:build integration
// +build integration

package v1

import (
	"context"
	"net/http"
	"testing"
	"time"

	customErrors "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/http_errors/custom_errors"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/test/messaging"
	"github.com/mehdihadeli/go-mediatr"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/suite"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/features/deleting_product/v1/commands"
	integrationEvents "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/features/deleting_product/v1/events/integration_events"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/shared/test_fixtures/integration"
)

type deleteProductIntegrationTests struct {
	*integration.IntegrationTestSharedFixture
}

func TestDeleteProductIntegration(t *testing.T) {
	suite.Run(
		t,
		&deleteProductIntegrationTests{
			IntegrationTestSharedFixture: integration.NewIntegrationTestSharedFixture(t),
		},
	)
}

func (c *deleteProductIntegrationTests) Test_Should_Delete_Product_From_DB() {
	ctx := context.Background()

	id := c.Items[0].ProductId
	command, err := commands.NewDeleteProduct(id)
	c.Require().NoError(err)

	result, err := mediatr.Send[*commands.DeleteProduct, *mediatr.Unit](ctx, command)

	c.Require().NoError(err)
	c.Assert().NotNil(result)

	deletedProduct, err := c.ProductRepository.GetProductById(ctx, id)
	c.Assert().Nil(deletedProduct)
}

func (c *deleteProductIntegrationTests) Test_Should_Returns_NotFound_Error_When_Record_DoesNot_Exists() {
	ctx := context.Background()

	id := uuid.NewV4()
	command, err := commands.NewDeleteProduct(id)
	c.Require().NoError(err)

	result, err := mediatr.Send[*commands.DeleteProduct, *mediatr.Unit](ctx, command)

	c.Assert().Error(err)
	c.True(customErrors.IsApplicationError(err, http.StatusNotFound))
	c.True(customErrors.IsNotFoundError(err))
	c.Assert().Nil(result)
}

func (c *deleteProductIntegrationTests) Test_Should_Publish_Product_Created_To_Broker() {
	ctx := context.Background()

	shouldPublish := messaging.ShouldProduced[*integrationEvents.ProductDeletedV1](
		ctx,
		c.Bus,
		nil,
	)

	id := c.Items[0].ProductId
	command, err := commands.NewDeleteProduct(id)
	c.Require().NoError(err)

	_, err = mediatr.Send[*commands.DeleteProduct, *mediatr.Unit](ctx, command)
	c.Require().NoError(err)

	// ensuring message published to the rabbitmq broker
	shouldPublish.Validate(ctx, "there is no published message", time.Second*30)
}
