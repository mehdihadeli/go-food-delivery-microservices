package v1

import (
	"context"
	"fmt"
	"github.com/mehdihadeli/go-mediatr"
	customErrors "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/http_errors/custom_errors"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/tracing"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/contracts"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type DeleteProductCommandHandler struct {
	log             logger.Logger
	cfg             *config.Config
	mongoRepository contracts.ProductRepository
	redisRepository contracts.ProductCacheRepository
}

func NewDeleteProductCommandHandler(log logger.Logger, cfg *config.Config, repository contracts.ProductRepository, redisRepository contracts.ProductCacheRepository) *DeleteProductCommandHandler {
	return &DeleteProductCommandHandler{log: log, cfg: cfg, mongoRepository: repository, redisRepository: redisRepository}
}

func (c *DeleteProductCommandHandler) Handle(ctx context.Context, command *DeleteProductCommand) (*mediatr.Unit, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "DeleteProductCommandHandler.Handle")
	span.LogFields(log.String("ProductId", command.ProductID.String()))
	span.LogFields(log.Object("Command", command))
	defer span.Finish()

	if err := c.mongoRepository.DeleteProductByID(ctx, command.ProductID); err != nil {
		return nil, tracing.TraceWithErr(span, customErrors.NewApplicationErrorWrap(err, "[DeleteProductCommandHandler_Handle.DeleteProductByID] error in deleting product in the mongo repository"))
	}

	c.log.Infof("(product deleted) id: {%s}", command.ProductID)

	err := c.redisRepository.DeleteProduct(ctx, command.ProductID.String())
	if err != nil {
		return nil, tracing.TraceWithErr(span, customErrors.NewApplicationErrorWrap(err, "[DeleteProductCommandHandler_Handle.DeleteProduct] error in deleting product in the redis repository"))
	}

	c.log.Infow(fmt.Sprintf("[DeleteProductCommandHandler.Handle] product with id: {%s} deleted", command.ProductID), logger.Fields{"productId": command.ProductID})

	return &mediatr.Unit{}, nil
}
