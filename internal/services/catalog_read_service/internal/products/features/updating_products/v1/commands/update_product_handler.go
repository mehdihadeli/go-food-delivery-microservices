package commands

import (
	"context"
	"fmt"

	customErrors "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/http_errors/custom_errors"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/otel/tracing"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/otel/tracing/attribute"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogreadservice/internal/products/contracts/data"

	"github.com/mehdihadeli/go-mediatr"
	attribute2 "go.opentelemetry.io/otel/attribute"
)

type UpdateProductHandler struct {
	log             logger.Logger
	mongoRepository data.ProductRepository
	redisRepository data.ProductCacheRepository
	tracer          tracing.AppTracer
}

func NewUpdateProductHandler(
	log logger.Logger,
	mongoRepository data.ProductRepository,
	redisRepository data.ProductCacheRepository,
	tracer tracing.AppTracer,
) *UpdateProductHandler {
	return &UpdateProductHandler{
		log:             log,
		mongoRepository: mongoRepository,
		redisRepository: redisRepository,
		tracer:          tracer,
	}
}

func (c *UpdateProductHandler) Handle(
	ctx context.Context,
	command *UpdateProduct,
) (*mediatr.Unit, error) {
	ctx, span := c.tracer.Start(ctx, "UpdateProductHandler.Handle")
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
					"[UpdateProductHandler_Handle.GetProductById] error in fetching product with productId %s in the mongo repository",
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
					"[UpdateProductHandler_Handle.GetProductById] product with productId %s not found",
					command.ProductId,
				),
			),
		)
	}

	product.Price = command.Price
	product.Name = command.Name
	product.Description = command.Description
	product.UpdatedAt = command.UpdatedAt

	_, err = c.mongoRepository.UpdateProduct(ctx, product)
	if err != nil {
		return nil, tracing.TraceErrFromSpan(
			span,
			customErrors.NewApplicationErrorWrap(
				err,
				"[UpdateProductHandler_Handle.UpdateProduct] error in updating product in the mongo repository",
			),
		)
	}

	err = c.redisRepository.PutProduct(ctx, product.Id, product)
	if err != nil {
		return nil, tracing.TraceErrFromSpan(
			span,
			customErrors.NewApplicationErrorWrap(
				err,
				"[UpdateProductHandler_Handle.PutProduct] error in updating product in the redis repository",
			),
		)
	}

	c.log.Infow(
		fmt.Sprintf("[UpdateProductHandler.Handle] product with id: {%s} updated", product.Id),
		logger.Fields{"ProductId": command.ProductId, "Id": product.Id},
	)

	return &mediatr.Unit{}, nil
}
