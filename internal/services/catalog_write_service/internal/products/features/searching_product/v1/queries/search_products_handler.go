package queries

import (
	"context"

	customErrors "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/http_errors/custom_errors"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/otel/tracing"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/utils"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/contracts/data"
	dto "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/dto/v1"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/features/searching_product/v1/dtos"
)

type SearchProductsHandler struct {
	log    logger.Logger
	pgRepo data.ProductRepository
	tracer tracing.AppTracer
}

func NewSearchProductsHandler(
	log logger.Logger,
	pgRepo data.ProductRepository,
	tracer tracing.AppTracer,
) *SearchProductsHandler {
	return &SearchProductsHandler{log: log, pgRepo: pgRepo, tracer: tracer}
}

func (c *SearchProductsHandler) Handle(
	ctx context.Context,
	query *SearchProducts,
) (*dtos.SearchProductsResponseDto, error) {
	products, err := c.pgRepo.SearchProducts(
		ctx,
		query.SearchText,
		query.ListQuery,
	)
	if err != nil {
		return nil, customErrors.NewApplicationErrorWrap(
			err,
			"error in searching products in the repository",
		)
	}

	listResultDto, err := utils.ListResultToListResultDto[*dto.ProductDto](
		products,
	)
	if err != nil {
		return nil, customErrors.NewApplicationErrorWrap(
			err,
			"error in the mapping ListResultToListResultDto",
		)
	}

	c.log.Info("products fetched")

	return &dtos.SearchProductsResponseDto{Products: listResultDto}, nil
}
