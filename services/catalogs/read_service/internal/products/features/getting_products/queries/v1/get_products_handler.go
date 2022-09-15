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
	gettingProductsDto "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/features/getting_products/dtos"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type GetProductsHandler struct {
	log             logger.Logger
	cfg             *config.Config
	mongoRepository contracts.ProductRepository
}

func NewGetProductsHandler(log logger.Logger, cfg *config.Config, mongoRepository contracts.ProductRepository) *GetProductsHandler {
	return &GetProductsHandler{log: log, cfg: cfg, mongoRepository: mongoRepository}
}

func (c *GetProductsHandler) Handle(ctx context.Context, query *GetProducts) (*gettingProductsDto.GetProductsResponseDto, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GetProductsHandler.Handle")
	span.LogFields(log.Object("Query", query))
	defer span.Finish()

	products, err := c.mongoRepository.GetAllProducts(ctx, query.ListQuery)
	if err != nil {
		return nil, tracing.TraceWithErr(span, customErrors.NewApplicationErrorWrap(err, "[GetProductsHandler_Handle.GetAllProducts] error in getting products in the repository"))
	}

	listResultDto, err := utils.ListResultToListResultDto[*dto.ProductDto](products)
	if err != nil {
		return nil, tracing.TraceWithErr(span, customErrors.NewApplicationErrorWrap(err, "[GetProductsHandler_Handle.ListResultToListResultDto] error in the mapping ListResultToListResultDto"))
	}

	c.log.Info("[GetProductsHandler.Handle] products fetched")

	return &gettingProductsDto.GetProductsResponseDto{Products: listResultDto}, nil
}
