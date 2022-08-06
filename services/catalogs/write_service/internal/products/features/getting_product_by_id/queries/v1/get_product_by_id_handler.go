package v1

import (
	"context"
	"fmt"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/http_errors"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/mapper"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/contracts"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/dto"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/getting_product_by_id/dtos"
	"github.com/opentracing/opentracing-go"
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
	defer span.Finish()

	product, err := q.pgRepo.GetProductById(ctx, query.ProductID)

	if err != nil {
		return nil, httpErrors.NewNotFoundError(fmt.Sprintf("product with id %s not found", query.ProductID))
	}

	productDto, err := mapper.Map[*dto.ProductDto](product)
	if err != nil {
		return nil, err
	}

	return &dtos.GetProductByIdResponseDto{Product: productDto}, nil
}
