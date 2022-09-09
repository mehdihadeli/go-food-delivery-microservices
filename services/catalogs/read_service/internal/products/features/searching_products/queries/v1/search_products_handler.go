package v1

import (
	"context"
	customErrors "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/http_errors/custom_errors"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/tracing"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/utils"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/contracts"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/dto"
	searchingProductsDtos "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/features/searching_products/dtos"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type SearchProductsQueryHandler struct {
	log             logger.Logger
	cfg             *config.Config
	mongoRepository contracts.ProductRepository
}

func NewSearchProductsQueryHandler(log logger.Logger, cfg *config.Config, repository contracts.ProductRepository) *SearchProductsQueryHandler {
	return &SearchProductsQueryHandler{log: log, cfg: cfg, mongoRepository: repository}
}

func (c *SearchProductsQueryHandler) Handle(ctx context.Context, query *SearchProductsQuery) (*searchingProductsDtos.SearchProductsResponseDto, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SearchProductsQueryHandler.Handle")
	span.LogFields(log.Object("Query", query))
	defer span.Finish()

	products, err := c.mongoRepository.SearchProducts(ctx, query.SearchText, query.ListQuery)
	if err != nil {
		return nil, tracing.TraceWithErr(span, customErrors.NewApplicationErrorWrap(err, "[SearchProductsQueryHandler_Handle.SearchProducts] error in searching products in the repository"))
	}

	listResultDto, err := utils.ListResultToListResultDto[*dto.ProductDto](products)
	if err != nil {
		return nil, tracing.TraceWithErr(span, customErrors.NewApplicationErrorWrap(err, "[GetProductsQueryHandler_Handle.ListResultToListResultDto] error in the mapping ListResultToListResultDto"))
	}
	c.log.Info("[SearchProductsQueryHandler.Handle] products fetched")

	return &searchingProductsDtos.SearchProductsResponseDto{Products: listResultDto}, nil
}
