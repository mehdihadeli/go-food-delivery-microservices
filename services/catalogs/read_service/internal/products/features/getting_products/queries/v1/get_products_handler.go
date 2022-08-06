package v1

import (
	"context"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/utils"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/contracts"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/dto"
	gettingProductsDto "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/features/getting_products/dtos"
	"github.com/opentracing/opentracing-go"
)

type GetProductsQueryHandler struct {
	log    logger.Logger
	cfg    *config.Config
	pgRepo contracts.ProductRepository
}

func NewGetProductsQueryHandler(log logger.Logger, cfg *config.Config, pgRepo contracts.ProductRepository) *GetProductsQueryHandler {
	return &GetProductsQueryHandler{log: log, cfg: cfg, pgRepo: pgRepo}
}

func (c *GetProductsQueryHandler) Handle(ctx context.Context, query *GetProductsQuery) (*gettingProductsDto.GetProductsResponseDto, error) {

	span, ctx := opentracing.StartSpanFromContext(ctx, "GetProductsQueryHandler.Handle")
	defer span.Finish()

	products, err := c.pgRepo.GetAllProducts(ctx, query.ListQuery)
	if err != nil {
		return nil, err
	}

	listResultDto, err := utils.ListResultToListResultDto[*dto.ProductDto](products)
	if err != nil {
		return nil, err
	}

	return &gettingProductsDto.GetProductsResponseDto{Products: listResultDto}, nil
}
