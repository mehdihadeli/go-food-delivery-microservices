package commands

import (
	"context"
	"fmt"

	"github.com/mehdihadeli/go-mediatr"
	attribute2 "go.opentelemetry.io/otel/attribute"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogreadservice/internal/products/contracts/data"

	customErrors "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/http_errors/custom_errors"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/otel/tracing"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/otel/tracing/attribute"
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
	ctx, span := c.tracer.Start(ctx, "DeleteProductCommand.Handle")
	span.SetAttributes(attribute2.String("ProductId", command.ProductId.String()))
	span.SetAttributes(attribute.Object("Command", command))
	defer span.End()

	product, err := c.mongoRepository.GetProductByProductId(ctx, command.ProductId.String())
	if err != nil {
		return nil, tracing.TraceErrFromSpan(
			span,
			customErrors.NewApplicationErrorWrap(
				err,
				fmt.Sprintf(
					"[DeleteProductHandler_Handle.GetProductById] error in fetching product with productId %s in the mongo repository",
					command.ProductId,
				),
			),
		)
	}
	if product == nil {
		return nil, tracing.TraceErrFromSpan(
			span,
			customErrors.NewNotFoundErrorWrap(
				err,
				fmt.Sprintf(
					"[DeleteProductHandler_Handle.GetProductById] product with productId %s not found",
					command.ProductId,
				),
			),
		)
	}

	if err := c.mongoRepository.DeleteProductByID(ctx, product.Id); err != nil {
		return nil, tracing.TraceErrFromSpan(
			span,
			customErrors.NewApplicationErrorWrap(
				err,
				"[DeleteProductHandler_Handle.DeleteProductByID] error in deleting product in the mongo repository",
			),
		)
	}

	c.log.Infof("(product deleted) id: {%s}", product.Id)

	err = c.redisRepository.DeleteProduct(ctx, product.Id)
	if err != nil {
		return nil, tracing.TraceErrFromSpan(
			span,
			customErrors.NewApplicationErrorWrap(
				err,
				"[DeleteProductHandler_Handle.DeleteProduct] error in deleting product in the redis repository",
			),
		)
	}

	c.log.Infow(
		fmt.Sprintf("[DeleteProductCommand.Handle] product with id: {%s} deleted", product.Id),
		logger.Fields{"ProductId": command.ProductId, "Id": product.Id},
	)

	return &mediatr.Unit{}, nil
}
