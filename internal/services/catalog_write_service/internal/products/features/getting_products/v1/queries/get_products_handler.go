package queries

import (
	"context"

	customErrors "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/http_errors/custom_errors"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/otel/tracing"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/otel/tracing/attribute"
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
	ctx, span := c.tracer.Start(ctx, "GetProductsHandler.Handle")
	span.SetAttributes(attribute.Object("Query", query))
	defer span.End()

	products, err := c.pgRepo.GetAllProducts(ctx, query.ListQuery)
	if err != nil {
		return nil, tracing.TraceErrFromSpan(
			span,
			customErrors.NewApplicationErrorWrap(
				err,
				"[GetProductsHandler_Handle.GetAllProducts] error in getting products in the repository",
			),
		)
	}

	listResultDto, err := utils.ListResultToListResultDto[*dto.ProductDto](products)
	if err != nil {
		return nil, tracing.TraceErrFromSpan(
			span,
			customErrors.NewApplicationErrorWrap(
				err,
				"[GetProductsHandler_Handle.ListResultToListResultDto] error in the mapping ListResultToListResultDto",
			),
		)
	}

	c.log.Info("[GetProductsHandler.Handle] products fetched")

	return &dtos.GetProductsResponseDto{Products: listResultDto}, nil
}
