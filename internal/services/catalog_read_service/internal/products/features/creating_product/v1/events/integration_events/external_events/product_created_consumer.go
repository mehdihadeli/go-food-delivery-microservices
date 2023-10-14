package externalEvents

import (
	"context"
	"fmt"

	customErrors "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/http_errors/custom_errors"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/messaging/consumer"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/messaging/types"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/otel/tracing"
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
			"command validation failed",
		)

		return validationErr
	}
	_, err = mediatr.Send[*commands.CreateProduct, *dtos.CreateProductResponseDto](
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
