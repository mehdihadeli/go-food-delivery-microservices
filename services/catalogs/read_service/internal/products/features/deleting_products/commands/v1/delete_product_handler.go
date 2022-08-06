package v1

import (
	"context"
	"github.com/mehdihadeli/go-mediatr"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/contracts"
	"github.com/opentracing/opentracing-go"
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
	defer span.Finish()

	if err := c.mongoRepository.DeleteProductByID(ctx, command.ProductID); err != nil {
		return nil, err
	}

	c.log.Infof("(product deleted) id: {%s}", command.ProductID)

	c.redisRepository.DelProduct(ctx, command.ProductID.String())

	return &mediatr.Unit{}, nil
}
