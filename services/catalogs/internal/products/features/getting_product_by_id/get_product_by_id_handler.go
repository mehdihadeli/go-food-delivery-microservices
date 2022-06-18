package getting_product_by_id

import (
	"context"
	"github.com/eyazici90/go-mediator/mediator"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/errors"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/products/infrastructure/repositories"
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

func (q *GetProductByIdHandler) Handle(ctx context.Context, msg mediator.Message) (*models.Product, error) {
	query, ok := msg.(GetProductById)
	if err := errors.CheckType(ok); err != nil {
		return nil, err
	}

	return q.pgRepo.GetProductById(ctx, query.ProductID)
}
