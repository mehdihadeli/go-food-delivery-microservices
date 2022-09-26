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
	deletingProductV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/features/deleting_products/commands/v1"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	uuid "github.com/satori/go.uuid"
)

type productDeletedConsumer struct {
	*delivery.ProductConsumersBase
}

func NewProductDeletedConsumer(productConsumerBase *delivery.ProductConsumersBase) *productDeletedConsumer {
	return &productDeletedConsumer{productConsumerBase}
}

func (c *productDeletedConsumer) Handle(ctx context.Context, consumeContext types2.MessageConsumeContextT[*ProductDeletedV1]) error {
	if consumeContext.Message() == nil {
		return nil
	}
	span, ctx := opentracing.StartSpanFromContext(ctx, "productDeletedConsumer.Handle")
	span.LogFields(log.Object("Message", consumeContext.Created()))
	defer span.Finish()

	productUUID, err := uuid.FromString(consumeContext.Message().ProductId)
	if err != nil {
		badRequestErr := customErrors.NewBadRequestErrorWrap(err, "[productDeletedConsumer_Handle.uuid.FromString] error in the converting uuid")
		c.Log.Errorf(fmt.Sprintf("[productDeletedConsumer_Handle.uuid.FromString] err: %v", tracing.TraceWithErr(span, badRequestErr)))
		c.CommitErrMessage()

		return err
	}

	command := deletingProductV1.NewDeleteProduct(productUUID)
	if err := c.Validator.StructCtx(ctx, command); err != nil {
		validationErr := customErrors.NewValidationErrorWrap(err, "[productDeletedConsumer_Handle.StructCtx] command validation failed")
		c.Log.Errorf(fmt.Sprintf("[productDeletedConsumer_Consume.StructCtx] err: {%v}", tracing.TraceWithErr(span, validationErr)))
		c.CommitErrMessage()

		return err
	}

	_, err = mediatr.Send[*deletingProductV1.DeleteProduct, *mediatr.Unit](ctx, command)

	if err != nil {
		err = errors.WithMessage(err, "[productDeletedConsumer_Handle.Send] error in sending DeleteProduct")
		c.Log.Errorw(fmt.Sprintf("[productDeletedConsumer_Handle.Send] id: {%s}, err: {%v}", command.ProductId, tracing.TraceWithErr(span, err)), logger.Fields{"Id": command.ProductId})

		c.CommitErrMessage()
	}
	c.CommitMessage()

	return nil
}
