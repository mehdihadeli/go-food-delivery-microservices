package v1

import (
	"context"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/utils"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/contracts"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/dto"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/searching_product/dtos"
	"github.com/opentracing/opentracing-go"
)

type SearchProductsQueryHandler struct {
	log    logger.Logger
	cfg    *config.Config
	pgRepo contracts.ProductRepository
}

func NewSearchProductsQueryHandler(log logger.Logger, cfg *config.Config, pgRepo contracts.ProductRepository) *SearchProductsQueryHandler {
	return &SearchProductsQueryHandler{log: log, cfg: cfg, pgRepo: pgRepo}
}

func (c *SearchProductsQueryHandler) Handle(ctx context.Context, query *SearchProductsQuery) (*dtos.SearchProductsResponseDto, error) {

	span, ctx := opentracing.StartSpanFromContext(ctx, "SearchProductsQueryHandler.Handle")
	defer span.Finish()

	products, err := c.pgRepo.SearchProducts(ctx, query.SearchText, query.ListQuery)
	if err != nil {
		return nil, err
	}

	listResultDto, err := utils.ListResultToListResultDto[*dto.ProductDto](products)
	if err != nil {
		return nil, err
	}

	return &dtos.SearchProductsResponseDto{Products: listResultDto}, nil
}
