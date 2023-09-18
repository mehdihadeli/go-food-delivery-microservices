package commands

import (
	"context"
	"fmt"
	"net/http"

	customErrors "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/http_errors/custom_errors"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/mapper"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/messaging/producer"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/otel/tracing"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/otel/tracing/attribute"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/contracts/data"
	dto "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/dto/v1"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/features/updating_product/v1/events/integration_events"

	"github.com/mehdihadeli/go-mediatr"
	attribute2 "go.opentelemetry.io/otel/attribute"
)

type UpdateProductHandler struct {
	log              logger.Logger
	uow              data.CatalogUnitOfWork
	rabbitmqProducer producer.Producer
	tracer           tracing.AppTracer
}

func NewUpdateProductHandler(
	log logger.Logger,
	uow data.CatalogUnitOfWork,
	rabbitmqProducer producer.Producer,
	tracer tracing.AppTracer,
) *UpdateProductHandler {
	return &UpdateProductHandler{
		log:              log,
		uow:              uow,
		rabbitmqProducer: rabbitmqProducer,
		tracer:           tracer,
	}
}

func (c *UpdateProductHandler) Handle(
	ctx context.Context,
	command *UpdateProduct,
) (*mediatr.Unit, error) {
	ctx, span := c.tracer.Start(ctx, "UpdateProductHandler.Handle")
	span.SetAttributes(attribute2.String("ProductId", command.ProductID.String()))
	span.SetAttributes(attribute.Object("Command", command))
	defer span.End()

	err := c.uow.Do(ctx, func(catalogContext data.CatalogContext) error {
		product, err := catalogContext.Products().GetProductById(ctx, command.ProductID)
		if err != nil {
			return tracing.TraceErrFromSpan(
				span,
				customErrors.NewApplicationErrorWrapWithCode(
					err,
					http.StatusNotFound,
					fmt.Sprintf(
						"[UpdateProductHandler_Handle.GetProductById] product with id %s not found",
						command.ProductID,
					),
				),
			)
		}

		product.Name = command.Name
		product.Price = command.Price
		product.Description = command.Description
		product.UpdatedAt = command.UpdatedAt

		updatedProduct, err := catalogContext.Products().UpdateProduct(ctx, product)
		if err != nil {
			return tracing.TraceErrFromSpan(
				span,
				customErrors.NewApplicationErrorWrap(
					err,
					"[UpdateProductHandler_Handle.UpdateProduct] error in updating product in the repository",
				),
			)
		}

		productDto, err := mapper.Map[*dto.ProductDto](updatedProduct)
		if err != nil {
			return tracing.TraceErrFromSpan(
				span,
				customErrors.NewApplicationErrorWrap(
					err,
					"[UpdateProductHandler_Handle.Map] error in the mapping ProductDto",
				),
			)
		}

		productUpdated := integration_events.NewProductUpdatedV1(productDto)

		err = c.rabbitmqProducer.PublishMessage(ctx, productUpdated, nil)
		if err != nil {
			return tracing.TraceErrFromSpan(
				span,
				customErrors.NewApplicationErrorWrap(
					err,
					"[UpdateProductHandler_Handle.PublishMessage] error in publishing 'ProductUpdated' message",
				),
			)
		}

		c.log.Infow(
			fmt.Sprintf(
				"[UpdateProductHandler.Handle] product with id '%s' updated",
				command.ProductID,
			),
			logger.Fields{"ProductId": command.ProductID},
		)

		c.log.Infow(
			fmt.Sprintf(
				"[DeleteProductHandler.Handle] ProductUpdated message with messageId `%s` published to the rabbitmq broker",
				productUpdated.MessageId,
			),
			logger.Fields{"MessageId": productUpdated.MessageId},
		)

		return nil
	})

	return &mediatr.Unit{}, err
}
