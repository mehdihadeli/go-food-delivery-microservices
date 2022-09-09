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
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/models"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type UpdateProductCommandHandler struct {
	log             logger.Logger
	cfg             *config.Config
	mongoRepository contracts.ProductRepository
	redisRepository contracts.ProductCacheRepository
}

func NewUpdateProductCommandHandler(log logger.Logger, cfg *config.Config, mongoRepository contracts.ProductRepository, redisRepository contracts.ProductCacheRepository) *UpdateProductCommandHandler {
	return &UpdateProductCommandHandler{log: log, cfg: cfg, mongoRepository: mongoRepository, redisRepository: redisRepository}
}

func (c *UpdateProductCommandHandler) Handle(ctx context.Context, command *UpdateProductCommand) (*mediatr.Unit, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UpdateProductCommandHandler.Handle")
	span.LogFields(log.String("ProductId", command.ProductID.String()))
	span.LogFields(log.Object("Command", command))
	defer span.Finish()

	p, err := c.mongoRepository.GetProductById(ctx, command.ProductID)
	if err != nil {
		return nil, customErrors.NewApplicationErrorWrap(err, fmt.Sprintf("[UpdateProductCommandHandler_Handle.GetProductById] error in fetching product with id %s in the mongo repository", command.ProductID))
	}

	if p == nil {
		return nil, customErrors.NewNotFoundErrorWrap(err, fmt.Sprintf("[UpdateProductCommandHandler_Handle.GetProductById] product with id %s not found", command.ProductID))
	}

	product := &models.Product{ProductID: command.ProductID.String(), Name: command.Name, Description: command.Description, Price: command.Price, UpdatedAt: command.UpdatedAt}

	updatedProduct, err := c.mongoRepository.UpdateProduct(ctx, product)
	if err != nil {
		return nil, tracing.TraceWithErr(span, customErrors.NewApplicationErrorWrap(err, "[UpdateProductCommandHandler_Handle.UpdateProduct] error in updating product in the mongo repository"))
	}

	err = c.redisRepository.PutProduct(ctx, updatedProduct.ProductID, updatedProduct)
	if err != nil {
		return nil, tracing.TraceWithErr(span, customErrors.NewApplicationErrorWrap(err, "[UpdateProductCommandHandler_Handle.PutProduct] error in updating product in the redis repository"))
	}

	c.log.Infow(fmt.Sprintf("[UpdateProductCommandHandler.Handle] product with id: {%s} updated", command.ProductID), logger.Fields{"productId": command.ProductID})

	return &mediatr.Unit{}, nil
}
