package updating_products

import (
	"context"
	"fmt"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/http_errors"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/mediatr"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/contracts"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/models"
	"github.com/opentracing/opentracing-go"
)

type UpdateProductHandler struct {
	log        logger.Logger
	cfg        *config.Config
	repository contracts.ProductRepository
}

func NewUpdateProductHandler(log logger.Logger, cfg *config.Config, repository contracts.ProductRepository) *UpdateProductHandler {
	return &UpdateProductHandler{log: log, cfg: cfg, repository: repository}
}

func (c *UpdateProductHandler) Handle(ctx context.Context, command *UpdateProduct) (*mediatr.Unit, error) {

	span, ctx := opentracing.StartSpanFromContext(ctx, "UpdateProductHandler.Handle")
	defer span.Finish()

	_, err := c.repository.GetProductById(ctx, command.ProductID)

	if err != nil {
		return nil, http_errors.NewNotFoundError(fmt.Sprintf("product with id %s not found", command.ProductID))
	}

	product := &models.Product{ProductID: command.ProductID.String(), Name: command.Name, Description: command.Description, Price: command.Price, UpdatedAt: command.UpdatedAt}

	_, err = c.repository.UpdateProduct(ctx, product)
	if err != nil {
		return nil, err
	}

	c.log.Infof("(product updated) id: {%s}", command.ProductID)

	return &mediatr.Unit{}, nil
}
