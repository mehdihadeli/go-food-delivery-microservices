package v1

import (
	"context"
	"emperror.dev/errors"
	"fmt"
	"github.com/mehdihadeli/go-mediatr"
	customErrors "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/http_errors/custom_errors"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	types2 "github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/types"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/tracing"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/delivery"
	updatingProductV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/features/updating_products/commands/v1"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	uuid "github.com/satori/go.uuid"
)

type productUpdatedConsumer struct {
	*delivery.ProductConsumersBase
}

func NewProductUpdatedConsumer(productConsumerBase *delivery.ProductConsumersBase) *productUpdatedConsumer {
	return &productUpdatedConsumer{productConsumerBase}
}

func (c *productUpdatedConsumer) Handle(ctx context.Context, consumeContext types2.MessageConsumeContextT[*ProductUpdatedV1]) error {
	if consumeContext.Message() == nil {
		return nil
	}
	span, ctx := opentracing.StartSpanFromContext(ctx, "productUpdatedConsumer.Handle")
	span.LogFields(log.Object("Message", consumeContext.Created()))
	defer span.Finish()

	updatedProduct := consumeContext.Message()

	productUUID, err := uuid.FromString(updatedProduct.ProductId)
	if err != nil {
		c.Log.WarnMsg("uuid.FromString", err)
		badRequestErr := customErrors.NewBadRequestErrorWrap(err, "[updateProductConsumer_Consume.uuid.FromString] error in the converting uuid")
		c.Log.Errorf(fmt.Sprintf("[updateProductConsumer_Consume.uuid.FromString] err: %v", tracing.TraceWithErr(span, badRequestErr)))
		c.CommitErrMessage()

		return err
	}

	command := updatingProductV1.NewUpdateProduct(productUUID, updatedProduct.Name, updatedProduct.Description, updatedProduct.Price)
	if err := c.Validator.StructCtx(ctx, command); err != nil {
		validationErr := customErrors.NewValidationErrorWrap(err, "[updateProductConsumer_Consume.StructCtx] command validation failed")
		c.Log.Errorf(fmt.Sprintf("[updateProductConsumer_Consume.StructCtx] err: {%v}", tracing.TraceWithErr(span, validationErr)))
		c.CommitErrMessage()

		return err
	}

	_, err = mediatr.Send[*updatingProductV1.UpdateProduct, *mediatr.Unit](ctx, command)
	if err != nil {
		err = errors.WithMessage(err, "[updateProductConsumer_Consume.Send] error in sending UpdateProduct")
		c.Log.Errorw(fmt.Sprintf("[updateProductConsumer_Consume.Send] id: {%s}, err: {%v}", command.ProductId, tracing.TraceWithErr(span, err)), logger.Fields{"Id": command.ProductId})
		c.CommitErrMessage()
		return err
	}
	c.CommitMessage()

	return nil
}
