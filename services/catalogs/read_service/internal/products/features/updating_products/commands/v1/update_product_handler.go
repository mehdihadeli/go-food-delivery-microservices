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

type UpdateProductHandler struct {
	log             logger.Logger
	cfg             *config.Config
	mongoRepository contracts.ProductRepository
	redisRepository contracts.ProductCacheRepository
}

func NewUpdateProductHandler(log logger.Logger, cfg *config.Config, mongoRepository contracts.ProductRepository, redisRepository contracts.ProductCacheRepository) *UpdateProductHandler {
	return &UpdateProductHandler{log: log, cfg: cfg, mongoRepository: mongoRepository, redisRepository: redisRepository}
}

func (c *UpdateProductHandler) Handle(ctx context.Context, command *UpdateProduct) (*mediatr.Unit, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UpdateProductHandler.Handle")
	span.LogFields(log.String("ProductId", command.ProductId.String()))
	span.LogFields(log.Object("Command", command))
	defer span.Finish()

	product, err := c.mongoRepository.GetProductByProductId(ctx, command.ProductId)
	if err != nil {
		return nil, customErrors.NewApplicationErrorWrap(err, fmt.Sprintf("[UpdateProductHandler_Handle.GetProductById] error in fetching product with productId %s in the mongo repository", command.ProductId))
	}

	if product == nil {
		return nil, customErrors.NewNotFoundErrorWrap(err, fmt.Sprintf("[UpdateProductHandler_Handle.GetProductById] product with productId %s not found", command.ProductId))
	}

	product.Price = command.Price
	product.Name = command.Name
	product.Description = command.Description
	product.UpdatedAt = command.UpdatedAt

	_, err = c.mongoRepository.UpdateProduct(ctx, product)
	if err != nil {
		return nil, tracing.TraceWithErr(span, customErrors.NewApplicationErrorWrap(err, "[UpdateProductHandler_Handle.UpdateProduct] error in updating product in the mongo repository"))
	}

	err = c.redisRepository.PutProduct(ctx, product.Id, product)
	if err != nil {
		return nil, tracing.TraceWithErr(span, customErrors.NewApplicationErrorWrap(err, "[UpdateProductHandler_Handle.PutProduct] error in updating product in the redis repository"))
	}

	c.log.Infow(fmt.Sprintf("[UpdateProductHandler.Handle] product with id: {%s} updated", product.Id), logger.Fields{"ProductId": command.ProductId, "Id": product.Id})

	return &mediatr.Unit{}, nil
}
