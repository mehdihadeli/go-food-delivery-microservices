package v1

import (
	"context"
	"fmt"
	customErrors "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/http_errors/custom_errors"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/mapper"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/otel/tracing"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/otel/tracing/attribute"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/contracts/repositories"
	ordersDto "github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/dtos"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/features/getting_order_by_id/dtos"
	attribute2 "go.opentelemetry.io/otel/attribute"
)

type GetOrderByIdHandler struct {
	log                  logger.Logger
	cfg                  *config.Config
	orderMongoRepository repositories.OrderReadRepository
}

func NewGetOrderByIdHandler(log logger.Logger, cfg *config.Config, orderMongoRepository repositories.OrderReadRepository) *GetOrderByIdHandler {
	return &GetOrderByIdHandler{log: log, cfg: cfg, orderMongoRepository: orderMongoRepository}
}

func (q *GetOrderByIdHandler) Handle(ctx context.Context, query *GetOrderById) (*dtos.GetOrderByIdResponseDto, error) {
	ctx, span := tracing.Tracer.Start(ctx, "GetOrderByIdHandler.Handle")
	span.SetAttributes(attribute2.String("Id", query.Id.String()))
	span.SetAttributes(attribute.Object("Query", query))
	defer span.End()

	// get order by order-read id
	order, err := q.orderMongoRepository.GetOrderById(ctx, query.Id)
	if err != nil {
		return nil, tracing.TraceErrFromSpan(span, customErrors.NewApplicationErrorWrap(err, fmt.Sprintf("[GetOrderByIdHandler_Handle.GetProductById] error in getting order with id %s in the mongo repository", query.Id.String())))
	}

	if order == nil {
		// get order by order-write id
		order, err = q.orderMongoRepository.GetOrderByOrderId(ctx, query.Id)
		if err != nil {
			return nil, tracing.TraceErrFromSpan(span, customErrors.NewApplicationErrorWrap(err, fmt.Sprintf("[GetOrderByIdHandler_Handle.GetProductById] error in getting order with orderId %s in the mongo repository", query.Id.String())))
		}
	}

	orderDto, err := mapper.Map[*ordersDto.OrderReadDto](order)
	if err != nil {
		return nil, tracing.TraceErrFromSpan(span, customErrors.NewApplicationErrorWrap(err, "[GetOrderByIdHandler_Handle.Map] error in the mapping order"))
	}

	q.log.Infow(fmt.Sprintf("[GetOrderByIdHandler.Handle] order with id: {%s} fetched", query.Id.String()), logger.Fields{"Id": query.Id})

	return &dtos.GetOrderByIdResponseDto{Order: orderDto}, nil
}
