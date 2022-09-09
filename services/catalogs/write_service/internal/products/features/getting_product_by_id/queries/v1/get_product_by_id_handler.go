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

type GetProductByIdQueryHandler struct {
	log    logger.Logger
	cfg    *config.Config
	pgRepo contracts.ProductRepository
}

func NewGetProductByIdQueryHandler(log logger.Logger, cfg *config.Config, pgRepo contracts.ProductRepository) *GetProductByIdQueryHandler {
	return &GetProductByIdQueryHandler{log: log, cfg: cfg, pgRepo: pgRepo}
}

func (q *GetProductByIdQueryHandler) Handle(ctx context.Context, query *GetProductByIdQuery) (*dtos.GetProductByIdResponseDto, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GetProductByIdQueryHandler.Handle")
	span.LogFields(log.String("ProductId", query.ProductID.String()))
	span.LogFields(log.Object("Query", query))
	defer span.Finish()

	product, err := q.pgRepo.GetProductById(ctx, query.ProductID)
	if err != nil {
		return nil, tracing.TraceWithErr(span, customErrors.NewApplicationErrorWrap(err, fmt.Sprintf("[GetProductByIdQueryHandler_Handle.GetProductById] error in getting product with id %d in the repository", query.ProductID)))
	}

	productDto, err := mapper.Map[*dto.ProductDto](product)
	if err != nil {
		return nil, tracing.TraceWithErr(span, customErrors.NewApplicationErrorWrap(err, "[GetProductByIdQueryHandler_Handle.Map] error in the mapping product"))
	}

	q.log.Infow(fmt.Sprintf("[GetProductByIdQueryHandler.Handle] product with id: {%d} fetched", query.ProductID), logger.Fields{"ProductId": query.ProductID})

	return &dtos.GetProductByIdResponseDto{Product: productDto}, nil
}
