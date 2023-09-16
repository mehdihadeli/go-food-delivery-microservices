package externalEvents

import (
	"context"
	"fmt"

	customErrors "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/http_errors/custom_errors"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/messaging/consumer"
	messageTracing "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/messaging/otel/tracing"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/messaging/types"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/otel/tracing"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/otel/tracing/attribute"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogreadservice/internal/products/features/creating_product/v1/commands"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogreadservice/internal/products/features/creating_product/v1/dtos"

	"emperror.dev/errors"
	"github.com/go-playground/validator"
	"github.com/mehdihadeli/go-mediatr"
)

type productCreatedConsumer struct {
	logger    logger.Logger
	validator *validator.Validate
	tracer    tracing.AppTracer
}

func NewProductCreatedConsumer(
	logger logger.Logger,
	validator *validator.Validate,
	tracer tracing.AppTracer,
) consumer.ConsumerHandler {
	return &productCreatedConsumer{logger: logger, validator: validator, tracer: tracer}
}

func (c *productCreatedConsumer) Handle(
	ctx context.Context,
	consumeContext types.MessageConsumeContext,
) error {
	product, ok := consumeContext.Message().(*ProductCreatedV1)
	if !ok {
		return errors.New("error in casting message to ProductCreatedV1")
	}

	ctx, span := c.tracer.Start(ctx, "productCreatedConsumer.Handle")
	span.SetAttributes(attribute.Object("Message", consumeContext.Message()))
	defer span.End()

	command, err := commands.NewCreateProduct(
		product.ProductId,
		product.Name,
		product.Description,
		product.Price,
		product.CreatedAt,
	)
	if err != nil {
		validationErr := customErrors.NewValidationErrorWrap(
			err,
			"[productCreatedConsumer_Handle.StructCtx] command validation failed",
		)
		c.logger.Errorf(
			fmt.Sprintf(
				"[productCreatedConsumer_Handle.StructCtx] err: {%v}",
				messageTracing.TraceMessagingErrFromSpan(span, validationErr),
			),
		)

		return err
	}
	_, err = mediatr.Send[*commands.CreateProduct, *dtos.CreateProductResponseDto](ctx, command)
	if err != nil {
		err = errors.WithMessage(
			err,
			"[productCreatedConsumer_Handle.Send] error in sending CreateProduct",
		)
		c.logger.Errorw(
			fmt.Sprintf(
				"[productCreatedConsumer_Handle.Send] id: {%s}, err: {%v}",
				command.ProductId,
				messageTracing.TraceMessagingErrFromSpan(span, err),
			),
			logger.Fields{"Id": command.ProductId},
		)
	}
	c.logger.Info("Product consumer handled.")

	return err
}
