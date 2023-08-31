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

	customErrors "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/http_errors/custom_errors"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/test/messaging"
	createProductCommand "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/features/creating_product/v1/commands"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/features/creating_product/v1/dtos"
	integrationEvents "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/features/creating_product/v1/events/integration_events"
)

func (p *productFeaturesIntegrationTestsv1) Test_Create_Product() {
	ctx := context.Background()

	p.T().Run("Should_Create_New_Product_To_DB", func(t *testing.T) {
		command, err := createProductCommand.NewCreateProduct(
			gofakeit.Name(),
			gofakeit.AdjectiveDescriptive(),
			gofakeit.Price(150, 6000),
		)
		p.Require().NoError(err)

		result, err := mediatr.Send[*createProductCommand.CreateProduct, *dtos.CreateProductResponseDto](
			ctx,
			command,
		)
		p.Require().NoError(err)

		p.Assert().NotNil(result)
		p.Assert().Equal(command.ProductID, result.ProductID)

		createdProduct, err := p.ProductRepository.GetProductById(
			ctx,
			result.ProductID,
		)
		p.Require().NoError(err)
		p.Assert().NotNil(createdProduct)
	})

	p.T().Run("Should_Return_Error_For_Duplicate_Record", func(t *testing.T) {
		id := p.Items[0].ProductId

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
		p.Assert().Error(err)
		p.True(customErrors.IsApplicationError(err, http.StatusConflict))
		p.Assert().Nil(result)
	})

	p.T().Run("Should_Publish_Product_Created_To_Broker", func(t *testing.T) {
		shouldPublish := messaging.ShouldProduced[*integrationEvents.ProductCreatedV1](
			ctx,
			p.Bus,
			nil,
		)

		command, err := createProductCommand.NewCreateProduct(
			gofakeit.Name(),
			gofakeit.AdjectiveDescriptive(),
			gofakeit.Price(150, 6000),
		)
		p.Require().NoError(err)

		_, err = mediatr.Send[*createProductCommand.CreateProduct, *dtos.CreateProductResponseDto](
			ctx,
			command,
		)
		p.Require().NoError(err)

		// ensuring message published to the rabbitmq broker
		shouldPublish.Validate(ctx, "there is no published message", time.Second*30)
	})

	p.T().Run("Should_Consume_Product_Created_With_Existing_Consumer_From_Broker", func(t *testing.T) {
		// we setup this handler in `BeforeTest`
		// we don't have a consumer in this service, so we simulate one consumer
		// check for consuming `ProductCreatedV1` message with existing consumer
		hypothesis := messaging.ShouldConsume[*integrationEvents.ProductCreatedV1](ctx, p.Bus, nil)

		command, err := createProductCommand.NewCreateProduct(
			gofakeit.Name(),
			gofakeit.AdjectiveDescriptive(),
			gofakeit.Price(150, 6000),
		)
		p.Require().NoError(err)

		_, err = mediatr.Send[*createProductCommand.CreateProduct, *dtos.CreateProductResponseDto](
			ctx,
			command,
		)
		p.Require().NoError(err)

		// ensuring message can be consumed with a consumer
		hypothesis.Validate(ctx, "there is no consumed message", time.Second*30)
	})

	p.T().Run("Should_Consume_Product_Created_With_New_Consumer_From_Broker", func(t *testing.T) {
		defer p.Bus.Stop()

		// check for consuming `ProductCreatedV1` message, with a new consumer
		hypothesis, err := messaging.ShouldConsumeNewConsumer[*integrationEvents.ProductCreatedV1](
			p.Bus,
		)
		p.Require().NoError(err)

		// at first, we should add new consumer to rabbitmq bus then start the broker, because we can't add new consumer after start.
		// we should also turn off consumer in `BeforeTest` for this test
		p.Bus.Start(ctx)

		// wait for consumers ready to consume before publishing messages, preparation background workers takes a bit time (for preventing messages lost)
		time.Sleep(1 * time.Second)

		command, err := createProductCommand.NewCreateProduct(
			gofakeit.Name(),
			gofakeit.AdjectiveDescriptive(),
			gofakeit.Price(150, 6000),
		)
		p.Require().NoError(err)

		_, err = mediatr.Send[*createProductCommand.CreateProduct, *dtos.CreateProductResponseDto](
			ctx,
			command,
		)
		p.Require().NoError(err)

		// ensuring message can be consumed with a consumer
		hypothesis.Validate(ctx, "there is no consumed message", time.Second*30)
	})
}
