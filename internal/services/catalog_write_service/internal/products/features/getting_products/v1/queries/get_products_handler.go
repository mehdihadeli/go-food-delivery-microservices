package queries

import (
	"context"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/otel/tracing"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/utils"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/contracts/data"
	dto "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/dto/v1"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/features/getting_products/v1/dtos"
)

type GetProductsHandler struct {
	log    logger.Logger
	pgRepo data.ProductRepository
	tracer tracing.AppTracer
}

func NewGetProductsHandler(
	log logger.Logger,
	pgRepo data.ProductRepository,
	tracer tracing.AppTracer,
) *GetProductsHandler {
	return &GetProductsHandler{log: log, pgRepo: pgRepo, tracer: tracer}
}

func (c *GetProductsHandler) Handle(
	ctx context.Context,
	query *GetProducts,
) (*dtos.GetProductsResponseDto, error) {
	products, err := c.pgRepo.GetAllProducts(ctx, query.ListQuery)
	if err != nil {
		return nil, err
	}

	listResultDto, err := utils.ListResultToListResultDto[*dto.ProductDto](
		products,
	)
	if err != nil {
		return nil, err
	}

	c.log.Info("products fetched")

	return &dtos.GetProductsResponseDto{Products: listResultDto}, nil
}
