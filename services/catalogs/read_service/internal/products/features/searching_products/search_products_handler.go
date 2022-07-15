package searching_products

import (
	"context"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/utils"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/contracts"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/dto"
	searching_products_dtos "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/features/searching_products/dtos"
	"github.com/opentracing/opentracing-go"
)

type SearchProductsHandler struct {
	log        logger.Logger
	cfg        *config.Config
	repository contracts.ProductRepository
}

func NewSearchProductsHandler(log logger.Logger, cfg *config.Config, repository contracts.ProductRepository) *SearchProductsHandler {
	return &SearchProductsHandler{log: log, cfg: cfg, repository: repository}
}

func (c *SearchProductsHandler) Handle(ctx context.Context, query *SearchProducts) (*searching_products_dtos.SearchProductsResponseDto, error) {

	span, ctx := opentracing.StartSpanFromContext(ctx, "SearchProductsHandler.Handle")
	defer span.Finish()

	products, err := c.repository.SearchProducts(ctx, query.SearchText, query.ListQuery)
	if err != nil {
		return nil, err
	}

	listResultDto, err := utils.ListResultToListResultDto[*dto.ProductDto](products)
	if err != nil {
		return nil, err
	}

	return &searching_products_dtos.SearchProductsResponseDto{Products: listResultDto}, nil
}
