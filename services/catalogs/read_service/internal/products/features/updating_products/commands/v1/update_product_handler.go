package v1

import (
	"context"
	"fmt"
	"github.com/mehdihadeli/go-mediatr"
	httpErrors "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http_errors"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/contracts"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/models"
	"github.com/opentracing/opentracing-go"
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
	defer span.Finish()

	_, err := c.mongoRepository.GetProductById(ctx, command.ProductID)

	if err != nil {
		return nil, httpErrors.NewNotFoundError(err, fmt.Sprintf("product with id %s not found", command.ProductID))
	}

	product := &models.Product{ProductID: command.ProductID.String(), Name: command.Name, Description: command.Description, Price: command.Price, UpdatedAt: command.UpdatedAt}

	updatedProduct, err := c.mongoRepository.UpdateProduct(ctx, product)
	if err != nil {
		return nil, err
	}

	c.redisRepository.PutProduct(ctx, updatedProduct.ProductID, updatedProduct)

	c.log.Infof("(product updated) id: {%s}", command.ProductID)

	return &mediatr.Unit{}, nil
}
