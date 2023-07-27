package queries

import (
	"context"
	"fmt"

	attribute2 "go.opentelemetry.io/otel/attribute"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogreadservice/internal/products/contracts/data"

	customErrors "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/http_errors/custom_errors"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/mapper"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/otel/tracing"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/otel/tracing/attribute"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogreadservice/internal/products/dto"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogreadservice/internal/products/features/get_product_by_id/v1/dtos"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogreadservice/internal/products/models"
)

type GetProductByIdHandler struct {
	log             logger.Logger
	mongoRepository data.ProductRepository
	redisRepository data.ProductCacheRepository
	tracer          tracing.AppTracer
}

func NewGetProductByIdHandler(
	log logger.Logger,
	mongoRepository data.ProductRepository,
	redisRepository data.ProductCacheRepository,
	tracer tracing.AppTracer,
) *GetProductByIdHandler {
	return &GetProductByIdHandler{
		log:             log,
		mongoRepository: mongoRepository,
		redisRepository: redisRepository,
		tracer:          tracer,
	}
}

func (q *GetProductByIdHandler) Handle(
	ctx context.Context,
	query *GetProductById,
) (*dtos.GetProductByIdResponseDto, error) {
	ctx, span := q.tracer.Start(ctx, "getProductByIdHandler.Handle")
	span.SetAttributes(attribute.Object("Query", query))
	span.SetAttributes(attribute2.String("Id", query.Id.String()))
	defer span.End()

	var product *models.Product
	redisProduct, err := q.redisRepository.GetProductById(ctx, query.Id.String())
	if err != nil {
		return nil, tracing.TraceErrFromSpan(
			span,
			customErrors.NewApplicationErrorWrap(
				err,
				fmt.Sprintf(
					"[GetProductByIdHandler_Handle.GetProductById] error in getting product with id %d in the redis repository",
					query.Id,
				),
			),
		)
	}

	if redisProduct != nil {
		product = redisProduct
	} else {
		var mongoProduct *models.Product
		mongoProduct, err = q.mongoRepository.GetProductById(ctx, query.Id.String())
		if err != nil {
			return nil, tracing.TraceErrFromSpan(span, customErrors.NewApplicationErrorWrap(err, fmt.Sprintf("[GetProductByIdHandler_Handle.GetProductById] error in getting product with id %d in the mongo repository", query.Id)))
		}
		if mongoProduct == nil {
			mongoProduct, err = q.mongoRepository.GetProductByProductId(ctx, query.Id.String())
		}
		if mongoProduct == nil {
			return nil, nil
		}

		product = mongoProduct
		err = q.redisRepository.PutProduct(ctx, product.Id, product)
		if err != nil {
			return new(dtos.GetProductByIdResponseDto), tracing.TraceErrFromSpan(span, err)
		}
	}

	productDto, err := mapper.Map[*dto.ProductDto](product)
	if err != nil {
		return nil, tracing.TraceErrFromSpan(
			span,
			customErrors.NewApplicationErrorWrap(
				err,
				"[GetProductByIdHandler_Handle.Map] error in the mapping product",
			),
		)
	}

	q.log.Infow(
		fmt.Sprintf("[GetProductByIdHandler.Handle] product with id: {%s} fetched", query.Id),
		logger.Fields{"ProductId": product.ProductId, "Id": product.Id},
	)

	return &dtos.GetProductByIdResponseDto{Product: productDto}, nil
}
