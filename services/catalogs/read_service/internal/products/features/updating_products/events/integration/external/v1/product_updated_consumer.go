package v1

import (
	"context"
	"emperror.dev/errors"
	"fmt"
	"github.com/mehdihadeli/go-mediatr"
	customErrors "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/http_errors/custom_errors"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	messageTracing "github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/otel/tracing"
	types2 "github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/types"
	tracing "github.com/mehdihadeli/store-golang-microservice-sample/pkg/otel/tracing"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/otel/tracing/attribute"
	updatingProductV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/features/updating_products/commands/v1"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/shared/configurations/infrastructure"
	uuid "github.com/satori/go.uuid"
)

type productUpdatedConsumer struct {
	*infrastructure.InfrastructureConfigurations
}

func NewProductUpdatedConsumer(infra *infrastructure.InfrastructureConfigurations) *productUpdatedConsumer {
	return &productUpdatedConsumer{InfrastructureConfigurations: infra}
}

func (c *productUpdatedConsumer) Handle(ctx context.Context, consumeContext types2.MessageConsumeContext) error {
	message, ok := consumeContext.Message().(*ProductUpdatedV1)
	if !ok {
		return errors.New("error in casting message to ProductUpdatedV1")
	}

	ctx, span := tracing.Tracer.Start(ctx, "productUpdatedConsumer.Handle")
	span.SetAttributes(attribute.Object("Message", consumeContext.Message()))
	defer span.End()

	productUUID, err := uuid.FromString(message.ProductId)
	if err != nil {
		c.Log.WarnMsg("uuid.FromString", err)
		badRequestErr := customErrors.NewBadRequestErrorWrap(err, "[updateProductConsumer_Consume.uuid.FromString] error in the converting uuid")
		c.Log.Errorf(fmt.Sprintf("[updateProductConsumer_Consume.uuid.FromString] err: %v", messageTracing.TraceMessagingErrFromSpan(span, badRequestErr)))
		return err
	}

	command := updatingProductV1.NewUpdateProduct(productUUID, message.Name, message.Description, message.Price)
	if err := c.Validator.StructCtx(ctx, command); err != nil {
		validationErr := customErrors.NewValidationErrorWrap(err, "[updateProductConsumer_Consume.StructCtx] command validation failed")
		c.Log.Errorf(fmt.Sprintf("[updateProductConsumer_Consume.StructCtx] err: {%v}", messageTracing.TraceMessagingErrFromSpan(span, validationErr)))
		return err
	}

	_, err = mediatr.Send[*updatingProductV1.UpdateProduct, *mediatr.Unit](ctx, command)
	if err != nil {
		err = errors.WithMessage(err, "[updateProductConsumer_Consume.Send] error in sending UpdateProduct")
		c.Log.Errorw(fmt.Sprintf("[updateProductConsumer_Consume.Send] id: {%s}, err: {%v}", command.ProductId, messageTracing.TraceMessagingErrFromSpan(span, err)), logger.Fields{"Id": command.ProductId})
		return err
	}

	return nil
}
