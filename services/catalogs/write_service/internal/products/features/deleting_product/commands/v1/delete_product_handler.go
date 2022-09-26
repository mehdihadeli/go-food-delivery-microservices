package v1

import (
	"context"
	"fmt"
	"github.com/mehdihadeli/go-mediatr"
	customErrors "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/http_errors/custom_errors"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/producer"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/tracing"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/contracts"
	v1 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/deleting_product/events/integration/v1"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type DeleteProductHandler struct {
	log              logger.Logger
	cfg              *config.Config
	pgRepo           contracts.ProductRepository
	rabbitmqProducer producer.Producer
}

func NewDeleteProductHandler(log logger.Logger, cfg *config.Config, pgRepo contracts.ProductRepository, rabbitmqProducer producer.Producer) *DeleteProductHandler {
	return &DeleteProductHandler{log: log, cfg: cfg, pgRepo: pgRepo, rabbitmqProducer: rabbitmqProducer}
}

func (c *DeleteProductHandler) Handle(ctx context.Context, command *DeleteProduct) (*mediatr.Unit, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "deleteProductHandler.Handle")
	span.LogFields(log.String("ProductId", command.ProductID.String()))
	span.LogFields(log.Object("Command", command))
	defer span.Finish()

	if err := c.pgRepo.DeleteProductByID(ctx, command.ProductID); err != nil {
		return nil, tracing.TraceWithErr(span, customErrors.NewApplicationErrorWrap(err, "[DeleteProductHandler_Handle.DeleteProductByID] error in deleting product in the repository"))
	}

	productDeleted := v1.NewProductDeletedV1(command.ProductID.String())
	err := c.rabbitmqProducer.PublishMessage(ctx, productDeleted, nil)
	if err != nil {
		return nil, tracing.TraceWithErr(span, customErrors.NewApplicationErrorWrap(err, "[DeleteProductHandler_Handle.PublishMessage] error in publishing 'ProductDeleted' message"))
	}

	c.log.Infow(fmt.Sprintf("[DeleteProductHandler.Handle] ProductDeleted message with messageId '%s' published to the rabbitmq broker", productDeleted.MessageId), logger.Fields{"MessageId": productDeleted.MessageId})

	c.log.Infow(fmt.Sprintf("[DeleteProductHandler.Handle] product with id '%s' deleted", command.ProductID), logger.Fields{"ProductId": command.ProductID})

	return &mediatr.Unit{}, nil
}
