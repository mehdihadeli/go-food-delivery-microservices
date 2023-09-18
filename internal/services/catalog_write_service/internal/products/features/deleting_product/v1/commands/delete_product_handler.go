package commands

import (
	"context"
	"fmt"
	"net/http"

	customErrors "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/http_errors/custom_errors"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/messaging/producer"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/otel/tracing"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/otel/tracing/attribute"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/contracts/data"
	integrationEvents "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/features/deleting_product/v1/events/integration_events"

	"github.com/mehdihadeli/go-mediatr"
	attribute2 "go.opentelemetry.io/otel/attribute"
)

type DeleteProductHandler struct {
	log              logger.Logger
	uow              data.CatalogUnitOfWork
	rabbitmqProducer producer.Producer
	tracer           tracing.AppTracer
}

func NewDeleteProductHandler(
	log logger.Logger,
	uow data.CatalogUnitOfWork,
	rabbitmqProducer producer.Producer,
	tracer tracing.AppTracer,
) *DeleteProductHandler {
	return &DeleteProductHandler{
		log:              log,
		uow:              uow,
		rabbitmqProducer: rabbitmqProducer,
		tracer:           tracer,
	}
}

func (c *DeleteProductHandler) Handle(
	ctx context.Context,
	command *DeleteProduct,
) (*mediatr.Unit, error) {
	ctx, span := c.tracer.Start(ctx, "deleteProductHandler.Handle")
	span.SetAttributes(attribute2.String("ProductId", command.ProductID.String()))
	span.SetAttributes(attribute.Object("Command", command))
	defer span.End()

	err := c.uow.Do(ctx, func(catalogContext data.CatalogContext) error {
		if err := catalogContext.Products().DeleteProductByID(ctx, command.ProductID); err != nil {
			return tracing.TraceErrFromSpan(
				span,
				customErrors.NewApplicationErrorWrapWithCode(
					err,
					http.StatusNotFound,
					"[DeleteProductHandler_Handle.DeleteProductByID] product not found",
				),
			)
		}

		productDeleted := integrationEvents.NewProductDeletedV1(command.ProductID.String())
		err := c.rabbitmqProducer.PublishMessage(ctx, productDeleted, nil)
		if err != nil {
			return tracing.TraceErrFromSpan(
				span,
				customErrors.NewApplicationErrorWrap(
					err,
					"[DeleteProductHandler_Handle.PublishMessage] error in publishing 'ProductDeleted' message",
				),
			)
		}

		c.log.Infow(
			fmt.Sprintf(
				"[DeleteProductHandler.Handle] ProductDeleted message with messageId '%s' published to the rabbitmq broker",
				productDeleted.MessageId,
			),
			logger.Fields{"MessageId": productDeleted.MessageId},
		)

		c.log.Infow(
			fmt.Sprintf(
				"[DeleteProductHandler.Handle] product with id '%s' deleted",
				command.ProductID,
			),
			logger.Fields{"ProductId": command.ProductID},
		)

		return nil
	})

	return &mediatr.Unit{}, err
}
