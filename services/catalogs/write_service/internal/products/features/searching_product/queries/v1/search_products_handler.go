package v1

import (
	"context"
	customErrors "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/http_errors/custom_errors"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/tracing"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/utils"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/contracts"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/dto"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/searching_product/dtos"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type SearchProductsHandler struct {
	log    logger.Logger
	cfg    *config.Config
	pgRepo contracts.ProductRepository
}

func NewSearchProductsHandler(log logger.Logger, cfg *config.Config, pgRepo contracts.ProductRepository) *SearchProductsHandler {
	return &SearchProductsHandler{log: log, cfg: cfg, pgRepo: pgRepo}
}

func (c *SearchProductsHandler) Handle(ctx context.Context, query *SearchProducts) (*dtos.SearchProductsResponseDto, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SearchProductsHandler.Handle")
	span.LogFields(log.Object("Query", query))
	defer span.Finish()

	products, err := c.pgRepo.SearchProducts(ctx, query.SearchText, query.ListQuery)
	if err != nil {
		return nil, tracing.TraceWithErr(span, customErrors.NewApplicationErrorWrap(err, "[SearchProductsHandler_Handle.SearchProducts] error in searching products in the repository"))
	}

	listResultDto, err := utils.ListResultToListResultDto[*dto.ProductDto](products)
	if err != nil {
		return nil, tracing.TraceWithErr(span, customErrors.NewApplicationErrorWrap(err, "[SearchProductsHandler_Handle.ListResultToListResultDto] error in the mapping ListResultToListResultDto"))
	}

	c.log.Info("[SearchProductsHandler.Handle] products fetched")

	return &dtos.SearchProductsResponseDto{Products: listResultDto}, nil
}
