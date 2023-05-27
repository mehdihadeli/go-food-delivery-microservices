package externalEvents

import (
	"context"
	"fmt"

	"emperror.dev/errors"
	"github.com/mehdihadeli/go-mediatr"
	uuid "github.com/satori/go.uuid"

	customErrors "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/http_errors/custom_errors"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
	messageTracing "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/messaging/otel/tracing"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/messaging/types"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/otel/tracing"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/otel/tracing/attribute"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/read_service/internal/products/features/deleting_products/v1/commands"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/read_service/internal/shared/contracts"
)

type productDeletedConsumer struct {
	*contracts.InfrastructureConfigurations
}

func NewProductDeletedConsumer(infra *contracts.InfrastructureConfigurations) *productDeletedConsumer {
	return &productDeletedConsumer{InfrastructureConfigurations: infra}
}

func (c *productDeletedConsumer) Handle(ctx context.Context, consumeContext types.MessageConsumeContext) error {
	message, ok := consumeContext.Message().(*ProductDeletedV1)
	if !ok {
		return errors.New("error in casting message to ProductDeletedV1")
	}

	ctx, span := tracing.Tracer.Start(ctx, "productDeletedConsumer.Handle")
	span.SetAttributes(attribute.Object("Message", consumeContext.Message()))
	defer span.End()

	productUUID, err := uuid.FromString(message.ProductId)
	if err != nil {
		badRequestErr := customErrors.NewBadRequestErrorWrap(err, "[productDeletedConsumer_Handle.uuid.FromString] error in the converting uuid")
		c.Log.Errorf(fmt.Sprintf("[productDeletedConsumer_Handle.uuid.FromString] err: %v", messageTracing.TraceMessagingErrFromSpan(span, badRequestErr)))

		return err
	}

	command := commands.NewDeleteProduct(productUUID)
	if err := c.Validator.StructCtx(ctx, command); err != nil {
		validationErr := customErrors.NewValidationErrorWrap(err, "[productDeletedConsumer_Handle.StructCtx] command validation failed")
		c.Log.Errorf(fmt.Sprintf("[productDeletedConsumer_Consume.StructCtx] err: {%v}", messageTracing.TraceMessagingErrFromSpan(span, validationErr)))

		return err
	}

	_, err = mediatr.Send[*commands.DeleteProduct, *mediatr.Unit](ctx, command)

	if err != nil {
		err = errors.WithMessage(err, "[productDeletedConsumer_Handle.Send] error in sending DeleteProduct")
		c.Log.Errorw(fmt.Sprintf("[productDeletedConsumer_Handle.Send] id: {%s}, err: {%v}", command.ProductId, messageTracing.TraceMessagingErrFromSpan(span, err)), logger.Fields{"Id": command.ProductId})
	}

	return nil
}
