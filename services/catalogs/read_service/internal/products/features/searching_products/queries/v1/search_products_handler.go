package v1

import (
	"context"
	customErrors "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/http_errors/custom_errors"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/otel/tracing"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/otel/tracing/attribute"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/utils"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/contracts"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/dto"
	searchingProductsDtos "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/features/searching_products/dtos"
)

type SearchProductsHandler struct {
	log             logger.Logger
	cfg             *config.Config
	mongoRepository contracts.ProductRepository
}

func NewSearchProductsHandler(log logger.Logger, cfg *config.Config, repository contracts.ProductRepository) *SearchProductsHandler {
	return &SearchProductsHandler{log: log, cfg: cfg, mongoRepository: repository}
}

func (c *SearchProductsHandler) Handle(ctx context.Context, query *SearchProducts) (*searchingProductsDtos.SearchProductsResponseDto, error) {
	ctx, span := tracing.Tracer.Start(ctx, "SearchProductsHandler.Handle")
	span.SetAttributes(attribute.Object("Query", query))
	defer span.End()

	products, err := c.mongoRepository.SearchProducts(ctx, query.SearchText, query.ListQuery)
	if err != nil {
		return nil, tracing.TraceErrFromSpan(span, customErrors.NewApplicationErrorWrap(err, "[SearchProductsHandler_Handle.SearchProducts] error in searching products in the repository"))
	}

	listResultDto, err := utils.ListResultToListResultDto[*dto.ProductDto](products)
	if err != nil {
		return nil, tracing.TraceErrFromSpan(span, customErrors.NewApplicationErrorWrap(err, "[SearchProductsHandler_Handle.ListResultToListResultDto] error in the mapping ListResultToListResultDto"))
	}
	c.log.Info("[SearchProductsHandler.Handle] products fetched")

	return &searchingProductsDtos.SearchProductsResponseDto{Products: listResultDto}, nil
}
