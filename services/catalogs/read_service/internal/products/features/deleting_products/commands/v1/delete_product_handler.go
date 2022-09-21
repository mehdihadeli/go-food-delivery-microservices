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

type DeleteProductCommand struct {
	log             logger.Logger
	cfg             *config.Config
	mongoRepository contracts.ProductRepository
	redisRepository contracts.ProductCacheRepository
}

func NewDeleteProductCommand(log logger.Logger, cfg *config.Config, repository contracts.ProductRepository, redisRepository contracts.ProductCacheRepository) *DeleteProductCommand {
	return &DeleteProductCommand{log: log, cfg: cfg, mongoRepository: repository, redisRepository: redisRepository}
}

func (c *DeleteProductCommand) Handle(ctx context.Context, command *DeleteProduct) (*mediatr.Unit, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "DeleteProductCommand.Handle")
	span.LogFields(log.String("ProductId", command.ProductId.String()))
	span.LogFields(log.Object("Command", command))
	defer span.Finish()

	if err := c.mongoRepository.DeleteProductByID(ctx, command.ProductId); err != nil {
		return nil, tracing.TraceWithErr(span, customErrors.NewApplicationErrorWrap(err, "[DeleteProductHandler_Handle.DeleteProductByID] error in deleting product in the mongo repository"))
	}

	c.log.Infof("(product deleted) id: {%s}", command.ProductId)

	err := c.redisRepository.DeleteProduct(ctx, command.ProductId.String())
	if err != nil {
		return nil, tracing.TraceWithErr(span, customErrors.NewApplicationErrorWrap(err, "[DeleteProductHandler_Handle.DeleteProduct] error in deleting product in the redis repository"))
	}

	c.log.Infow(fmt.Sprintf("[DeleteProductCommand.Handle] product with id: {%s} deleted", command.ProductId), logger.Fields{"productId": command.ProductId})

	return &mediatr.Unit{}, nil
}
