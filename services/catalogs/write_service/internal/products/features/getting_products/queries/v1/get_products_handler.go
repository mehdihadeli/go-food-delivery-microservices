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
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/getting_products/dtos"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type GetProductsQueryHandler struct {
	log    logger.Logger
	cfg    *config.Config
	pgRepo contracts.ProductRepository
}

func NewGetProductsQueryHandler(log logger.Logger, cfg *config.Config, pgRepo contracts.ProductRepository) *GetProductsQueryHandler {
	return &GetProductsQueryHandler{log: log, cfg: cfg, pgRepo: pgRepo}
}

func (c *GetProductsQueryHandler) Handle(ctx context.Context, query *GetProductsQuery) (*dtos.GetProductsResponseDto, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GetProductsQueryHandler.Handle")
	span.LogFields(log.Object("Query", query))
	defer span.Finish()

	products, err := c.pgRepo.GetAllProducts(ctx, query.ListQuery)
	if err != nil {
		return nil, tracing.TraceWithErr(span, customErrors.NewApplicationErrorWrap(err, "[GetProductsQueryHandler_Handle.GetAllProducts] error in getting products in the repository"))
	}

	listResultDto, err := utils.ListResultToListResultDto[*dto.ProductDto](products)
	if err != nil {
		return nil, tracing.TraceWithErr(span, customErrors.NewApplicationErrorWrap(err, "[GetProductsQueryHandler_Handle.ListResultToListResultDto] error in the mapping ListResultToListResultDto"))
	}

	c.log.Info("[GetProductsQueryHandler.Handle] products fetched")

	return &dtos.GetProductsResponseDto{Products: listResultDto}, nil
}
