package getting_product_by_id

import (
	"context"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/products/contracts/repositories"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/products/dto"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/products/mappers"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/products/models"
)

type GetProductByIdQrHandler interface {
	Handle(ctx context.Context, query *GetProductById) (*models.Product, error)
}

type GetProductByIdHandler struct {
	log    logger.Logger
	cfg    *config.Config
	pgRepo repositories.ProductRepository
}

func NewGetProductByIdHandler(log logger.Logger, cfg *config.Config, pgRepo repositories.ProductRepository) *GetProductByIdHandler {
	return &GetProductByIdHandler{log: log, cfg: cfg, pgRepo: pgRepo}
}

func (q *GetProductByIdHandler) Handle(ctx context.Context, query GetProductById) (*dto.GetProductResponseDto, error) {
	product, err := q.pgRepo.GetProductById(ctx, query.ProductID)

	if err != nil {
		return nil, err
	}

	return mappers.ProductToGetProductResponseDto(product), nil
}
