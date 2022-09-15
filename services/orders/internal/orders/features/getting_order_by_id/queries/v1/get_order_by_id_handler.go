package v1

import (
	"context"
	"fmt"
	customErrors "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/http_errors/custom_errors"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/mapper"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/tracing"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/contracts/repositories"
	ordersDto "github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/dtos"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/features/getting_order_by_id/dtos"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type GetOrderByIdHandler struct {
	log                  logger.Logger
	cfg                  *config.Config
	orderMongoRepository repositories.OrderReadRepository
}

// TODO: Should read from read side model (mongo)

func NewGetOrderByIdHandler(log logger.Logger, cfg *config.Config, orderMongoRepository repositories.OrderReadRepository) *GetOrderByIdHandler {
	return &GetOrderByIdHandler{log: log, cfg: cfg, orderMongoRepository: orderMongoRepository}
}

func (q *GetOrderByIdHandler) Handle(ctx context.Context, query *GetOrderByIdQuery) (*dtos.GetOrderByIdResponseDto, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GetOrderByIdHandler.Handle")
	span.LogFields(log.String("ProductId", query.OrderId.String()))
	span.LogFields(log.Object("Query", query))
	defer span.Finish()

	order, err := q.orderMongoRepository.GetOrderById(ctx, query.OrderId)
	if err != nil {
		return nil, tracing.TraceWithErr(span, customErrors.NewApplicationErrorWrap(err, fmt.Sprintf("[GetOrderByIdHandler_Handle.GetProductById] error in getting order with id %s in the mongo repository", query.OrderId.String())))
	}

	orderDto, err := mapper.Map[*ordersDto.OrderReadDto](order)
	if err != nil {
		return nil, tracing.TraceWithErr(span, customErrors.NewApplicationErrorWrap(err, "[GetOrderByIdHandler_Handle.Map] error in the mapping order"))
	}

	q.log.Infow(fmt.Sprintf("[GetOrderByIdHandler.Handle] order with id: {%s} fetched", query.OrderId.String()), logger.Fields{"OrderId": query.OrderId})

	return &dtos.GetOrderByIdResponseDto{Order: orderDto}, nil
}
