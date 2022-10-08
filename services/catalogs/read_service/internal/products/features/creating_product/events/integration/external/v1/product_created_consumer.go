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
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/otel/tracing"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/otel/tracing/attribute"
	creatingProduct "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/features/creating_product"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/features/creating_product/commands/v1"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/shared/contracts"
)

type productCreatedConsumer struct {
	contracts.InfrastructureConfigurations
}

func NewProductCreatedConsumer(infra contracts.InfrastructureConfigurations) *productCreatedConsumer {
	return &productCreatedConsumer{InfrastructureConfigurations: infra}
}

func (c *productCreatedConsumer) Handle(ctx context.Context, consumeContext types2.MessageConsumeContext) error {
	product, ok := consumeContext.Message().(*ProductCreatedV1)
	if !ok {
		return errors.New("error in casting message to ProductCreatedV1")
	}

	ctx, span := tracing.Tracer.Start(ctx, "productCreatedConsumer.Handle")
	span.SetAttributes(attribute.Object("Message", consumeContext.Message()))
	defer span.End()

	command := v1.NewCreateProduct(product.ProductId, product.Name, product.Description, product.Price, product.CreatedAt)
	if err := c.Validator().StructCtx(ctx, command); err != nil {
		validationErr := customErrors.NewValidationErrorWrap(err, "[productCreatedConsumer_Handle.StructCtx] command validation failed")
		c.Log().Errorf(fmt.Sprintf("[productCreatedConsumer_Handle.StructCtx] err: {%v}", messageTracing.TraceMessagingErrFromSpan(span, validationErr)))

		return err
	}
	_, err := mediatr.Send[*v1.CreateProduct, *creatingProduct.CreateProductResponseDto](ctx, command)

	if err != nil {
		err = errors.WithMessage(err, "[productCreatedConsumer_Handle.Send] error in sending CreateProduct")
		c.Log().Errorw(fmt.Sprintf("[productCreatedConsumer_Handle.Send] id: {%s}, err: {%v}", command.ProductId, messageTracing.TraceMessagingErrFromSpan(span, err)), logger.Fields{"Id": command.ProductId})
	}

	return nil
}
