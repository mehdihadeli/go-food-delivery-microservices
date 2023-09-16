package queries

import (
	"context"
	"fmt"

	customErrors "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/http_errors/custom_errors"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/mapper"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/otel/tracing"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/otel/tracing/attribute"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orderservice/internal/orders/contracts/repositories"
	dtosV1 "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orderservice/internal/orders/dtos/v1"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orderservice/internal/orders/features/getting_order_by_id/v1/dtos"

	attribute2 "go.opentelemetry.io/otel/attribute"
)

type GetOrderByIdHandler struct {
	log                  logger.Logger
	orderMongoRepository repositories.OrderMongoRepository
	tracer               tracing.AppTracer
}

func NewGetOrderByIdHandler(
	log logger.Logger,
	orderMongoRepository repositories.OrderMongoRepository,
	tracer tracing.AppTracer,
) *GetOrderByIdHandler {
	return &GetOrderByIdHandler{
		log:                  log,
		orderMongoRepository: orderMongoRepository,
		tracer:               tracer,
	}
}

func (q *GetOrderByIdHandler) Handle(
	ctx context.Context,
	query *GetOrderById,
) (*dtos.GetOrderByIdResponseDto, error) {
	ctx, span := q.tracer.Start(ctx, "GetOrderByIdHandler.Handle")
	span.SetAttributes(attribute2.String("Id", query.Id.String()))
	span.SetAttributes(attribute.Object("Query", query))
	defer span.End()

	// get order by order-read id
	order, err := q.orderMongoRepository.GetOrderById(ctx, query.Id)
	if err != nil {
		return nil, tracing.TraceErrFromSpan(
			span,
			customErrors.NewApplicationErrorWrap(
				err,
				fmt.Sprintf(
					"[GetOrderByIdHandler_Handle.GetProductById] error in getting order with id %s in the mongo repository",
					query.Id.String(),
				),
			),
		)
	}

	if order == nil {
		// get order by order-write id
		order, err = q.orderMongoRepository.GetOrderByOrderId(ctx, query.Id)
		if err != nil {
			return nil, tracing.TraceErrFromSpan(
				span,
				customErrors.NewApplicationErrorWrap(
					err,
					fmt.Sprintf(
						"[GetOrderByIdHandler_Handle.GetProductById] error in getting order with orderId %s in the mongo repository",
						query.Id.String(),
					),
				),
			)
		}
	}

	orderDto, err := mapper.Map[*dtosV1.OrderReadDto](order)
	if err != nil {
		return nil, tracing.TraceErrFromSpan(
			span,
			customErrors.NewApplicationErrorWrap(
				err,
				"[GetOrderByIdHandler_Handle.Map] error in the mapping order",
			),
		)
	}

	q.log.Infow(
		fmt.Sprintf("[GetOrderByIdHandler.Handle] order with id: {%s} fetched", query.Id.String()),
		logger.Fields{"Id": query.Id},
	)

	return &dtos.GetOrderByIdResponseDto{Order: orderDto}, nil
}
