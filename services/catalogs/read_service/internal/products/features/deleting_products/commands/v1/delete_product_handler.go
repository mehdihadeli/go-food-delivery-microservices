package v1

import (
	"context"
	"fmt"
	"github.com/mehdihadeli/go-mediatr"
	customErrors "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/http_errors/custom_errors"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/otel/tracing"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/otel/tracing/attribute"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/contracts"
	uuid "github.com/satori/go.uuid"
	attribute2 "go.opentelemetry.io/otel/attribute"
)

type DeleteProductCommand struct {
	log             logger.Logger
	cfg             *config.Config
	mongoRepository contracts.ProductRepository
	redisRepository contracts.ProductCacheRepository
}

func NewDeleteProductHandler(log logger.Logger, cfg *config.Config, repository contracts.ProductRepository, redisRepository contracts.ProductCacheRepository) *DeleteProductCommand {
	return &DeleteProductCommand{log: log, cfg: cfg, mongoRepository: repository, redisRepository: redisRepository}
}

func (c *DeleteProductCommand) Handle(ctx context.Context, command *DeleteProduct) (*mediatr.Unit, error) {
	ctx, span := tracing.Tracer.Start(ctx, "DeleteProductCommand.Handle")
	span.SetAttributes(attribute2.String("ProductId", command.ProductId.String()))
	span.SetAttributes(attribute.Object("Command", command))
	defer span.End()

	product, err := c.mongoRepository.GetProductByProductId(ctx, command.ProductId)
	if err != nil {
		return nil, tracing.TraceErrFromSpan(span, customErrors.NewApplicationErrorWrap(err, fmt.Sprintf("[DeleteProductHandler_Handle.GetProductById] error in fetching product with productId %s in the mongo repository", command.ProductId)))
	}
	if product == nil {
		return nil, tracing.TraceErrFromSpan(span, customErrors.NewNotFoundErrorWrap(err, fmt.Sprintf("[DeleteProductHandler_Handle.GetProductById] product with productId %s not found", command.ProductId)))
	}

	id, err := uuid.FromString(product.Id)
	if err != nil {
		return nil, tracing.TraceErrFromSpan(span, err)
	}

	if err := c.mongoRepository.DeleteProductByID(ctx, id); err != nil {
		return nil, tracing.TraceErrFromSpan(span, customErrors.NewApplicationErrorWrap(err, "[DeleteProductHandler_Handle.DeleteProductByID] error in deleting product in the mongo repository"))
	}

	c.log.Infof("(product deleted) id: {%s}", id.String())

	err = c.redisRepository.DeleteProduct(ctx, product.Id)
	if err != nil {
		return nil, tracing.TraceErrFromSpan(span, customErrors.NewApplicationErrorWrap(err, "[DeleteProductHandler_Handle.DeleteProduct] error in deleting product in the redis repository"))
	}

	c.log.Infow(fmt.Sprintf("[DeleteProductCommand.Handle] product with id: {%s} deleted", id.String()), logger.Fields{"ProductId": command.ProductId, "Id": id})

	return &mediatr.Unit{}, nil
}
