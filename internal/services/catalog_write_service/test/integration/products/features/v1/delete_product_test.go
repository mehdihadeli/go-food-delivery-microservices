//go:build integration
// +build integration

package v1

import (
    "context"
    "net/http"
    "testing"

    "github.com/mehdihadeli/go-mediatr"
    uuid "github.com/satori/go.uuid"

    customErrors "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/http_errors/custom_errors"
    "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/features/deleting_product/v1/commands"
)

func (p *productFeaturesIntegrationTestsv1) Test_Delete_Product() {
	ctx := context.Background()

	p.T().Run("Should_Delete_Product_From_DB", func(t *testing.T) {
		id := p.Items[0].ProductId
		command, err := commands.NewDeleteProduct(id)
		p.Require().NoError(err)

		result, err := mediatr.Send[*commands.DeleteProduct, *mediatr.Unit](ctx, command)

		p.Require().NoError(err)
		p.Assert().NotNil(result)

		deletedProduct, err := p.ProductRepository.GetProductById(ctx, id)
		p.Assert().Nil(deletedProduct)
	})

	p.T().Run("Should_Returns_NotFound_Error_When_Record_DoesNot_Exists", func(t *testing.T) {
		id := uuid.NewV4()
		command, err := commands.NewDeleteProduct(id)
		p.Require().NoError(err)

		result, err := mediatr.Send[*commands.DeleteProduct, *mediatr.Unit](ctx, command)

		p.Assert().Error(err)
		p.True(customErrors.IsApplicationError(err, http.StatusNotFound))
		p.True(customErrors.IsNotFoundError(err))
		p.Assert().Nil(result)
	})

	//p.T().Run("Should_Publish_Product_Deleted_To_Broker", func(t *testing.T) {
	//	shouldPublish := messaging.ShouldProduced[*integrationEvents.ProductDeletedV1](
	//		ctx,
	//		p.Bus,
	//		nil,
	//	)
	//
	//	id := p.Items[0].ProductId
	//	command, err := commands.NewDeleteProduct(id)
	//	p.Require().NoError(err)
	//
	//	_, err = mediatr.Send[*commands.DeleteProduct, *mediatr.Unit](ctx, command)
	//	p.Require().NoError(err)
	//
	//	// ensuring message published to the rabbitmq broker
	//	shouldPublish.Validate(ctx, "there is no published message", time.Second*30)
	//})
}
