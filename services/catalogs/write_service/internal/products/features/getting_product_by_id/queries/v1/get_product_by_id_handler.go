package v1

import (
	"context"
	"fmt"
	customErrors "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/http_errors/custom_errors"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/mapper"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/tracing"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/contracts"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/dto"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/getting_product_by_id/dtos"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type GetProductByIdHandler struct {
	log    logger.Logger
	cfg    *config.Config
	pgRepo contracts.ProductRepository
}

func NewGetProductByIdHandler(log logger.Logger, cfg *config.Config, pgRepo contracts.ProductRepository) *GetProductByIdHandler {
	return &GetProductByIdHandler{log: log, cfg: cfg, pgRepo: pgRepo}
}

func (q *GetProductByIdHandler) Handle(ctx context.Context, query *GetProductById) (*dtos.GetProductByIdResponseDto, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GetProductByIdHandler.Handle")
	span.LogFields(log.String("ProductId", query.ProductID.String()))
	span.LogFields(log.Object("Query", query))
	defer span.Finish()

	product, err := q.pgRepo.GetProductById(ctx, query.ProductID)
	if err != nil {
		return nil, tracing.TraceWithErr(span, customErrors.NewApplicationErrorWrap(err, fmt.Sprintf("[GetProductByIdHandler_Handle.GetProductById] error in getting product with id %d in the repository", query.ProductID)))
	}

	productDto, err := mapper.Map[*dto.ProductDto](product)
	if err != nil {
		return nil, tracing.TraceWithErr(span, customErrors.NewApplicationErrorWrap(err, "[GetProductByIdHandler_Handle.Map] error in the mapping product"))
	}

	q.log.Infow(fmt.Sprintf("[GetProductByIdHandler.Handle] product with id: {%d} fetched", query.ProductID), logger.Fields{"ProductId": query.ProductID})

	return &dtos.GetProductByIdResponseDto{Product: productDto}, nil
}
