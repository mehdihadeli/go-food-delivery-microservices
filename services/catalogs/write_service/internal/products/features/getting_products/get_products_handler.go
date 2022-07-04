package getting_products

import (
	"context"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/contracts"

	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/mappers"
	"github.com/opentracing/opentracing-go"
)

type GetProductsHandler struct {
	log    logger.Logger
	cfg    *config.Config
	pgRepo contracts.ProductRepository
}

func NewGetProductsHandler(log logger.Logger, cfg *config.Config, pgRepo contracts.ProductRepository) *GetProductsHandler {
	return &GetProductsHandler{log: log, cfg: cfg, pgRepo: pgRepo}
}

func (c *GetProductsHandler) Handle(ctx context.Context, query GetProducts) (*GetProductsResponseDto, error) {

	span, ctx := opentracing.StartSpanFromContext(ctx, "GetProductsHandler.Handle")
	defer span.Finish()

	products, err := c.pgRepo.GetAllProducts(ctx, query.ListQuery)
	if err != nil {
		return nil, err
	}

	listResultDto := mappers.ListResultToListResultDto(products, mappers.ProductsToProductsDto)

	return &GetProductsResponseDto{Products: listResultDto}, nil
}
