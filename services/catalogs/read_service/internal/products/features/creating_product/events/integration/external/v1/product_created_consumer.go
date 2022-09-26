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
	creatingProduct "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/features/creating_product"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/features/creating_product/commands/v1"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type productCreatedConsumer struct {
	*delivery.ProductConsumersBase
}

func NewProductCreatedConsumer(productConsumerBase *delivery.ProductConsumersBase) *productCreatedConsumer {
	return &productCreatedConsumer{productConsumerBase}
}

func (c *productCreatedConsumer) Handle(ctx context.Context, consumeContext types2.MessageConsumeContextT[*ProductCreatedV1]) error {
	if consumeContext.Message() == nil {
		return nil
	}
	span, ctx := opentracing.StartSpanFromContext(ctx, "productCreatedConsumer.Handle")
	span.LogFields(log.Object("Message", consumeContext.Created()))
	defer span.Finish()

	product := consumeContext.Message()

	command := v1.NewCreateProduct(product.ProductId, product.Name, product.Description, product.Price, product.CreatedAt)
	if err := c.Validator.StructCtx(ctx, command); err != nil {
		validationErr := customErrors.NewValidationErrorWrap(err, "[productCreatedConsumer_Handle.StructCtx] command validation failed")
		c.Log.Errorf(fmt.Sprintf("[productCreatedConsumer_Handle.StructCtx] err: {%v}", tracing.TraceWithErr(span, validationErr)))
		c.CommitErrMessage()

		return err
	}
	_, err := mediatr.Send[*v1.CreateProduct, *creatingProduct.CreateProductResponseDto](ctx, command)

	if err != nil {
		err = errors.WithMessage(err, "[productCreatedConsumer_Handle.Send] error in sending CreateProduct")
		c.Log.Errorw(fmt.Sprintf("[productCreatedConsumer_Handle.Send] id: {%s}, err: {%v}", command.ProductId, tracing.TraceWithErr(span, err)), logger.Fields{"Id": command.ProductId})
		c.CommitErrMessage()
	}
	c.CommitMessage()

	return nil
}
