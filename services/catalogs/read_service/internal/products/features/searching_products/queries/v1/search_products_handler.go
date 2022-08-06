package v1

import (
	"context"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/utils"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/contracts"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/dto"
	searchingProductsDtos "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/features/searching_products/dtos"
	"github.com/opentracing/opentracing-go"
)

type SearchProductsQueryHandler struct {
	log        logger.Logger
	cfg        *config.Config
	repository contracts.ProductRepository
}

func NewSearchProductsQueryHandler(log logger.Logger, cfg *config.Config, repository contracts.ProductRepository) *SearchProductsQueryHandler {
	return &SearchProductsQueryHandler{log: log, cfg: cfg, repository: repository}
}

func (c *SearchProductsQueryHandler) Handle(ctx context.Context, query *SearchProductsQuery) (*searchingProductsDtos.SearchProductsResponseDto, error) {

	span, ctx := opentracing.StartSpanFromContext(ctx, "SearchProductsQueryHandler.Handle")
	defer span.Finish()

	products, err := c.repository.SearchProducts(ctx, query.SearchText, query.ListQuery)
	if err != nil {
		return nil, err
	}

	listResultDto, err := utils.ListResultToListResultDto[*dto.ProductDto](products)
	if err != nil {
		return nil, err
	}

	return &searchingProductsDtos.SearchProductsResponseDto{Products: listResultDto}, nil
}
