package externalEvents

import (
	"context"
	"fmt"

	"emperror.dev/errors"
	"github.com/mehdihadeli/go-mediatr"

	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/http_errors/custom_errors"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	messageTracing "github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/otel/tracing"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/types"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/otel/tracing"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/otel/tracing/attribute"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/features/creating_product/v1/commands"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/features/creating_product/v1/dtos"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/shared/contracts"
)

type productCreatedConsumer struct {
	*contracts.InfrastructureConfigurations
}

func NewProductCreatedConsumer(infra *contracts.InfrastructureConfigurations) *productCreatedConsumer {
	return &productCreatedConsumer{InfrastructureConfigurations: infra}
}

func (c *productCreatedConsumer) Handle(ctx context.Context, consumeContext types.MessageConsumeContext) error {
	product, ok := consumeContext.Message().(*ProductCreatedV1)
	if !ok {
		return errors.New("error in casting message to ProductCreatedV1")
	}

	ctx, span := tracing.Tracer.Start(ctx, "productCreatedConsumer.Handle")
	span.SetAttributes(attribute.Object("Message", consumeContext.Message()))
	defer span.End()

	command, err := commands.NewCreateProduct(product.ProductId, product.Name, product.Description, product.Price, product.CreatedAt)
	if err != nil {
		validationErr := customErrors.NewValidationErrorWrap(err, "[productCreatedConsumer_Handle.StructCtx] command validation failed")
		c.Log.Errorf(fmt.Sprintf("[productCreatedConsumer_Handle.StructCtx] err: {%v}", messageTracing.TraceMessagingErrFromSpan(span, validationErr)))

		return err
	}
	_, err = mediatr.Send[*commands.CreateProduct, *dtos.CreateProductResponseDto](ctx, command)
	if err != nil {
		err = errors.WithMessage(err, "[productCreatedConsumer_Handle.Send] error in sending CreateProduct")
		c.Log.Errorw(fmt.Sprintf("[productCreatedConsumer_Handle.Send] id: {%s}, err: {%v}", command.ProductId, messageTracing.TraceMessagingErrFromSpan(span, err)), logger.Fields{"Id": command.ProductId})
	}

	return nil
}
