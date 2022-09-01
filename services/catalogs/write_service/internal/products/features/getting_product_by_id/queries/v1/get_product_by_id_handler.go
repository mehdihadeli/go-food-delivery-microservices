package v1

import (
	"context"
	"fmt"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/mapper"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/tracing"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/contracts"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/dto"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/getting_product_by_id/dtos"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
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
	span, ctx := opentracing.StartSpanFromContext(ctx, "getProductByIdHandler.Handle")
	span.LogFields(log.Object("Query", query))
	defer span.Finish()

	product, err := q.pgRepo.GetProductById(ctx, query.ProductID)

	if err != nil {
		return nil, tracing.TraceWithErr(span, errors.WithMessage(err, fmt.Sprintf("[GetProductByIdQueryHandler_Handle.GetProductById] error in getting product with id %d in the repository", query.ProductID)))
	}

	productDto, err := mapper.Map[*dto.ProductDto](product)
	if err != nil {
		return nil, tracing.TraceWithErr(span, errors.Wrap(err, "[GetProductByIdQueryHandler_Handle.Map] error in the mapping product"))
	}

	return &dtos.GetProductByIdResponseDto{Product: productDto}, nil
}
