package getting_product_by_id

import (
	"context"
	"fmt"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/http_errors"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/contracts"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/getting_product_by_id/dtos"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/mappers"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/models"
)

type GetProductByIdQrHandler interface {
	Handle(ctx context.Context, query *GetProductById) (*models.Product, error)
}

type GetProductByIdHandler struct {
	log    logger.Logger
	cfg    *config.Config
	pgRepo contracts.ProductRepository
}

func NewGetProductByIdHandler(log logger.Logger, cfg *config.Config, pgRepo contracts.ProductRepository) *GetProductByIdHandler {
	return &GetProductByIdHandler{log: log, cfg: cfg, pgRepo: pgRepo}
}

func (q *GetProductByIdHandler) Handle(ctx context.Context, query GetProductById) (*dtos.GetProductByIdResponseDto, error) {
	product, err := q.pgRepo.GetProductById(ctx, query.ProductID)

	if err != nil {
		return nil, http_errors.NewNotFoundError(fmt.Sprintf("product with id %s not found", query.ProductID))
	}

	productDto := mappers.ProductToProductDto(product)

	return &dtos.GetProductByIdResponseDto{Product: productDto}, nil
}
