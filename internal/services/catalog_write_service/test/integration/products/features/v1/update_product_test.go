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
    uuid "github.com/satori/go.uuid"
    "github.com/stretchr/testify/require"

    customErrors "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/http_errors/custom_errors"
    "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/test/messaging"
    "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/features/updating_product/v1/commands"
    "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/features/updating_product/v1/events/integration_events"
)

func (p *productFeaturesIntegrationTestsv1) Test_Update_Product() {
	ctx := context.Background()

	p.T().Run("Should_Update_Existing_Product_In_DB", func(t *testing.T) {
		existing := p.Items[0]

		command, err := commands.NewUpdateProduct(
			existing.ProductId,
			gofakeit.Name(),
			existing.Description,
			existing.Price,
		)
		p.Require().NoError(err)

		result, err := mediatr.Send[*commands.UpdateProduct, *mediatr.Unit](ctx, command)
		p.Require().NoError(err)

		p.NotNil(result)

		updatedProduct, err := p.ProductRepository.GetProductById(
			ctx,
			existing.ProductId,
		)
		p.NotNil(updatedProduct)
		p.Equal(existing.ProductId, updatedProduct.ProductId)
		p.Equal(existing.Price, updatedProduct.Price)
		p.NotEqual(existing.Name, updatedProduct.Name)
	})

	p.T().Run("Should_Return_NotFound_Error_When_Item_DoesNot_Exist", func(t *testing.T) {
		id := uuid.NewV4()

		command, err := commands.NewUpdateProduct(
			id,
			gofakeit.Name(),
			gofakeit.EmojiDescription(),
			gofakeit.Price(150, 6000),
		)
		p.Require().NoError(err)

		result, err := mediatr.Send[*commands.UpdateProduct, *mediatr.Unit](ctx, command)

		p.Assert().Error(err)
		p.True(customErrors.IsApplicationError(err, http.StatusNotFound))
		p.Assert().Nil(result)
	})

	p.T().Run("Should_Return_NotFound_Error_When_Item_DoesNot_Exist", func(t *testing.T) {
		shouldPublish := messaging.ShouldProduced[*integration_events.ProductUpdatedV1](
			ctx,
			p.Bus,
			nil,
		)

		existing := p.Items[0]

		command, err := commands.NewUpdateProduct(
			existing.ProductId,
			gofakeit.Name(),
			existing.Description,
			existing.Price,
		)
		p.Require().NoError(err)

		_, err = mediatr.Send[*commands.UpdateProduct, *mediatr.Unit](ctx, command)
		p.Require().NoError(err)

		// ensuring message published to the rabbitmq broker
		shouldPublish.Validate(ctx, "there is no published message", time.Second*30)
	})

	p.T().Run("Should_Consume_Product_Created_With_Existing_Consumer_From_Broker", func(t *testing.T) {
		// we don't have a consumer in this service, so we simulate one consumer in `SetupSuite`
		// // check for consuming `ProductUpdatedV1` message with existing consumer
		hypothesis := messaging.ShouldConsume[*integration_events.ProductUpdatedV1](ctx, p.Bus, nil)

		existing := p.Items[0]
		command, err := commands.NewUpdateProduct(
			existing.ProductId,
			gofakeit.Name(),
			existing.Description,
			existing.Price,
		)
		p.Require().NoError(err)

		_, err = mediatr.Send[*commands.UpdateProduct, *mediatr.Unit](ctx, command)
		p.Require().NoError(err)

		// ensuring message can be consumed with a consumer
		hypothesis.Validate(ctx, "there is no consumed message", time.Second*30)
	})

	p.T().Run("Should_Consume_Product_Updated_With_New_Consumer_From_Broker", func(t *testing.T) {
		//  check for consuming `ProductUpdatedV1` message, with a new consumer
		hypothesis, err := messaging.ShouldConsumeNewConsumer[*integration_events.ProductUpdatedV1](
			p.Bus,
		)
		require.NoError(p.T(), err)

		// at first, we should add new consumer to rabbitmq bus then start the broker, because we can't add new consumer after start.
		// we should also turn off consumer in `BeforeTest` for this test
		p.Bus.Start(ctx)

		// wait for consumers ready to consume before publishing messages, preparation background workers takes a bit time (for preventing messages lost)
		time.Sleep(1 * time.Second)

		existing := p.Items[0]
		command, err := commands.NewUpdateProduct(
			existing.ProductId,
			gofakeit.Name(),
			existing.Description,
			existing.Price,
		)
		p.Require().NoError(err)

		_, err = mediatr.Send[*commands.UpdateProduct, *mediatr.Unit](ctx, command)
		p.Require().NoError(err)

		// ensuring message can be consumed with a consumer
		hypothesis.Validate(ctx, "there is no consumed message", time.Second*30)
	})
}
