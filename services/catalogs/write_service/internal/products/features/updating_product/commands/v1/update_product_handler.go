package v1

import (
	"context"
	"fmt"
	"github.com/mehdihadeli/go-mediatr"
	customErrors "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/http_errors/custom_errors"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/mapper"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/producer"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/otel/tracing"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/otel/tracing/attribute"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/contracts"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/dto"
	v1 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/updating_product/events/integration/v1"
	attribute2 "go.opentelemetry.io/otel/attribute"
)

type UpdateProductHandler struct {
	log              logger.Logger
	cfg              *config.Config
	pgRepo           contracts.ProductRepository
	rabbitmqProducer producer.Producer
}

func NewUpdateProductHandler(log logger.Logger, cfg *config.Config, pgRepo contracts.ProductRepository, rabbitmqProducer producer.Producer) *UpdateProductHandler {
	return &UpdateProductHandler{log: log, cfg: cfg, pgRepo: pgRepo, rabbitmqProducer: rabbitmqProducer}
}

func (c *UpdateProductHandler) Handle(ctx context.Context, command *UpdateProduct) (*mediatr.Unit, error) {
	ctx, span := tracing.Tracer.Start(ctx, "UpdateProductHandler.Handle")
	span.SetAttributes(attribute2.String("ProductId", command.ProductID.String()))
	span.SetAttributes(attribute.Object("Command", command))
	defer span.End()

	product, err := c.pgRepo.GetProductById(ctx, command.ProductID)
	if err != nil {
		return nil, tracing.TraceErrFromSpan(span, customErrors.NewApplicationErrorWrap(err, fmt.Sprintf("[UpdateProductHandler_Handle.GetProductById] error in fetching product with id %s", command.ProductID)))
	}

	if product == nil {
		return nil, tracing.TraceErrFromSpan(span, customErrors.NewNotFoundErrorWrap(err, fmt.Sprintf("[UpdateProductHandler_Handle.GetProductById] product with id %s not found", command.ProductID)))
	}

	product.Name = command.Name
	product.Price = command.Price
	product.Description = command.Description
	product.UpdatedAt = command.UpdatedAt

	updatedProduct, err := c.pgRepo.UpdateProduct(ctx, product)
	if err != nil {
		return nil, tracing.TraceErrFromSpan(span, customErrors.NewApplicationErrorWrap(err, "[UpdateProductHandler_Handle.UpdateProduct] error in updating product in the repository"))
	}

	productDto, err := mapper.Map[*dto.ProductDto](updatedProduct)
	if err != nil {
		return nil, tracing.TraceErrFromSpan(span, customErrors.NewApplicationErrorWrap(err, "[UpdateProductHandler_Handle.Map] error in the mapping ProductDto"))
	}

	productUpdated := v1.NewProductUpdatedV1(productDto)

	err = c.rabbitmqProducer.PublishMessage(ctx, productUpdated, nil)
	if err != nil {
		return nil, tracing.TraceErrFromSpan(span, customErrors.NewApplicationErrorWrap(err, "[UpdateProductHandler_Handle.PublishMessage] error in publishing 'ProductUpdated' message"))
	}

	c.log.Infow(fmt.Sprintf("[UpdateProductHandler.Handle] product with id '%s' updated", command.ProductID), logger.Fields{"ProductId": command.ProductID})

	c.log.Infow(fmt.Sprintf("[DeleteProductHandler.Handle] ProductUpdated message with messageId `%s` published to the rabbitmq broker", productUpdated.MessageId), logger.Fields{"MessageId": productUpdated.MessageId})

	return &mediatr.Unit{}, nil
}
