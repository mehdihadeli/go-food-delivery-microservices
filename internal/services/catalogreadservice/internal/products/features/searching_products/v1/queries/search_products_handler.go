package queries

import (
	"context"

	customErrors "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/http/httperrors/customerrors"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/logger"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/otel/tracing"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/utils"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogreadservice/internal/products/contracts/data"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogreadservice/internal/products/dto"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogreadservice/internal/products/features/searching_products/v1/dtos"
)

type SearchProductsHandler struct {
	log             logger.Logger
	mongoRepository data.ProductRepository
	tracer          tracing.AppTracer
}

func NewSearchProductsHandler(
	log logger.Logger,
	repository data.ProductRepository,
	tracer tracing.AppTracer,
) *SearchProductsHandler {
	return &SearchProductsHandler{
		log:             log,
		mongoRepository: repository,
		tracer:          tracer,
	}
}

func (c *SearchProductsHandler) Handle(
	ctx context.Context,
	query *SearchProducts,
) (*dtos.SearchProductsResponseDto, error) {
	products, err := c.mongoRepository.SearchProducts(
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
