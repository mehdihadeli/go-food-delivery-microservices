package externalEvents

import (
	"context"
	"fmt"

	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/core/messaging/consumer"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/core/messaging/types"
	customErrors "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/http/httperrors/customerrors"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/logger"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/otel/tracing"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/otel/tracing/attribute"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/otel/tracing/utils"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogreadservice/internal/products/features/updating_products/v1/commands"

	"emperror.dev/errors"
	"github.com/go-playground/validator"
	"github.com/mehdihadeli/go-mediatr"
	uuid "github.com/satori/go.uuid"
)

type productUpdatedConsumer struct {
	logger    logger.Logger
	validator *validator.Validate
	tracer    tracing.AppTracer
}

func NewProductUpdatedConsumer(
	logger logger.Logger,
	validator *validator.Validate,
	tracer tracing.AppTracer,
) consumer.ConsumerHandler {
	return &productUpdatedConsumer{
		logger:    logger,
		validator: validator,
		tracer:    tracer,
	}
}

func (c *productUpdatedConsumer) Handle(
	ctx context.Context,
	consumeContext types.MessageConsumeContext,
) error {
	message, ok := consumeContext.Message().(*ProductUpdatedV1)
	if !ok {
		return errors.New("error in casting message to ProductUpdatedV1")
	}

	ctx, span := c.tracer.Start(ctx, "productUpdatedConsumer.Handle")
	span.SetAttributes(attribute.Object("Message", consumeContext.Message()))
	defer span.End()

	productUUID, err := uuid.FromString(message.ProductId)
	if err != nil {
		c.logger.WarnMsg("uuid.FromString", err)
		badRequestErr := customErrors.NewBadRequestErrorWrap(
			err,
			"[updateProductConsumer_Consume.uuid.FromString] error in the converting uuid",
		)
		c.logger.Errorf(
			fmt.Sprintf(
				"[updateProductConsumer_Consume.uuid.FromString] err: %v",
				utils.TraceErrStatusFromSpan(span, badRequestErr),
			),
		)
		return err
	}

	command, err := commands.NewUpdateProduct(
		productUUID,
		message.Name,
		message.Description,
		message.Price,
	)
	if err != nil {
		validationErr := customErrors.NewValidationErrorWrap(
			err,
			"[updateProductConsumer_Consume.NewValidationErrorWrap] command validation failed",
		)
		c.logger.Errorf(
			fmt.Sprintf(
				"[updateProductConsumer_Consume.StructCtx] err: {%v}",
				utils.TraceErrStatusFromSpan(span, validationErr),
			),
		)
		return err
	}

	_, err = mediatr.Send[*commands.UpdateProduct, *mediatr.Unit](ctx, command)
	if err != nil {
		err = errors.WithMessage(
			err,
			"[updateProductConsumer_Consume.Send] error in sending UpdateProduct",
		)
		c.logger.Errorw(
			fmt.Sprintf(
				"[updateProductConsumer_Consume.Send] id: {%s}, err: {%v}",
				command.ProductId,
				utils.TraceErrStatusFromSpan(span, err),
			),
			logger.Fields{"Id": command.ProductId},
		)
		return err
	}

	return nil
}
