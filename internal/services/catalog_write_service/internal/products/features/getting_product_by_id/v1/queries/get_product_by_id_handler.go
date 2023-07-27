package getProductByIdQuery

import (
	"context"
	"fmt"
	"net/http"

	attribute2 "go.opentelemetry.io/otel/attribute"

	customErrors "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/http_errors/custom_errors"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/mapper"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/otel/tracing"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/otel/tracing/attribute"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/contracts/data"
	dtoV1 "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/dto/v1"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/features/getting_product_by_id/v1/dtos"
)

type GetProductByIdHandler struct {
	log    logger.Logger
	pgRepo data.ProductRepository
	tracer tracing.AppTracer
}

func NewGetProductByIdHandler(
	log logger.Logger,
	pgRepo data.ProductRepository,
	tracer tracing.AppTracer,
) *GetProductByIdHandler {
	return &GetProductByIdHandler{log: log, pgRepo: pgRepo, tracer: tracer}
}

func (q *GetProductByIdHandler) Handle(
	ctx context.Context,
	query *GetProductById,
) (*dtos.GetProductByIdResponseDto, error) {
	ctx, span := q.tracer.Start(ctx, "GetProductByIdHandler.Handle")
	span.SetAttributes(attribute.Object("Query", query))
	span.SetAttributes(attribute2.String("ProductId", query.ProductID.String()))
	defer span.End()

	product, err := q.pgRepo.GetProductById(ctx, query.ProductID)
	if err != nil {
		return nil, tracing.TraceErrFromSpan(
			span,
			customErrors.NewApplicationErrorWrapWithCode(
				err,
				http.StatusNotFound,
				fmt.Sprintf(
					"[GetProductByIdHandler_Handle.GetProductById] error in getting product with id %s in the repository",
					query.ProductID.String(),
				),
			),
		)
	}

	productDto, err := mapper.Map[*dtoV1.ProductDto](product)
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
		fmt.Sprintf(
			"[GetProductByIdHandler.Handle] product with id: {%s} fetched",
			query.ProductID,
		),
		logger.Fields{"ProductId": query.ProductID.String()},
	)

	return &dtos.GetProductByIdResponseDto{Product: productDto}, nil
}
