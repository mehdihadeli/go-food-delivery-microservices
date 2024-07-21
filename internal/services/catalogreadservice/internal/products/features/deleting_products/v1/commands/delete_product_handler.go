package commands

import (
	"context"
	"fmt"

	customErrors "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/http/httperrors/customerrors"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/logger"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/otel/tracing"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogreadservice/internal/products/contracts/data"

	"github.com/mehdihadeli/go-mediatr"
)

type DeleteProductCommand struct {
	log             logger.Logger
	mongoRepository data.ProductRepository
	redisRepository data.ProductCacheRepository
	tracer          tracing.AppTracer
}

func NewDeleteProductHandler(
	log logger.Logger,
	repository data.ProductRepository,
	redisRepository data.ProductCacheRepository,
	tracer tracing.AppTracer,
) *DeleteProductCommand {
	return &DeleteProductCommand{
		log:             log,
		mongoRepository: repository,
		redisRepository: redisRepository,
		tracer:          tracer,
	}
}

func (c *DeleteProductCommand) Handle(
	ctx context.Context,
	command *DeleteProduct,
) (*mediatr.Unit, error) {
	product, err := c.mongoRepository.GetProductByProductId(
		ctx,
		command.ProductId.String(),
	)
	if err != nil {
		return nil, customErrors.NewApplicationErrorWrap(
			err,
			fmt.Sprintf(
				"error in fetching product with productId %s in the mongo repository",
				command.ProductId,
			),
		)
	}
	if product == nil {
		return nil, customErrors.NewNotFoundErrorWrap(
			err,
			fmt.Sprintf(
				"product with productId %s not found",
				command.ProductId,
			),
		)
	}

	if err := c.mongoRepository.DeleteProductByID(ctx, product.Id); err != nil {
		return nil, customErrors.NewApplicationErrorWrap(
			err,
			"error in deleting product in the mongo repository",
		)
	}

	err = c.redisRepository.DeleteProduct(ctx, product.Id)
	if err != nil {
		return nil, customErrors.NewApplicationErrorWrap(
			err,
			"error in deleting product in the redis repository",
		)
	}

	c.log.Infow(
		fmt.Sprintf(
			"product with id: {%s} deleted",
			product.Id,
		),
		logger.Fields{"ProductId": command.ProductId, "Id": product.Id},
	)

	return &mediatr.Unit{}, nil
}
