package externalEvents

import (
	"context"
	"fmt"

	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/core/messaging/consumer"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/core/messaging/types"
	customErrors "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/http/httperrors/customerrors"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/logger"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/otel/tracing"
	v1 "github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogreadservice/internal/products/features/creating_product/v1"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogreadservice/internal/products/features/creating_product/v1/dtos"

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
	return &productCreatedConsumer{
		logger:    logger,
		validator: validator,
		tracer:    tracer,
	}
}

func (c *productCreatedConsumer) Handle(
	ctx context.Context,
	consumeContext types.MessageConsumeContext,
) error {
	product, ok := consumeContext.Message().(*ProductCreatedV1)
	if !ok {
		return errors.New("error in casting message to ProductCreatedV1")
	}

	command, err := v1.NewCreateProduct(
		product.ProductId,
		product.Name,
		product.Description,
		product.Price,
		product.CreatedAt,
	)
	if err != nil {
		validationErr := customErrors.NewValidationErrorWrap(
			err,
			"command validation failed",
		)

		return validationErr
	}
	_, err = mediatr.Send[*v1.CreateProduct, *dtos.CreateProductResponseDto](
		ctx,
		command,
	)
	if err != nil {
		return errors.WithMessage(
			err,
			fmt.Sprintf(
				"error in sending CreateProduct with id: {%s}",
				command.ProductId,
			),
		)
	}
	c.logger.Info("Product consumer handled.")

	return err
}
