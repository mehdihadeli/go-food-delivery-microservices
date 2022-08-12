package v1

import (
	"context"
	"fmt"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/es/contracts/store"
	httpErrors "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http_errors"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/mapper"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/config"
	ordersDto "github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/dtos"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/features/getting_order_by_id/dtos"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/models/orders/aggregate"
	"github.com/opentracing/opentracing-go"
)

type GetOrderByIdHandler struct {
	log            logger.Logger
	cfg            *config.Config
	aggregateStore store.AggregateStore[*aggregate.Order]
}

func NewGetOrderByIdHandler(log logger.Logger, cfg *config.Config, aggregateStore store.AggregateStore[*aggregate.Order]) *GetOrderByIdHandler {
	return &GetOrderByIdHandler{log: log, cfg: cfg, aggregateStore: aggregateStore}
}

func (q *GetOrderByIdHandler) Handle(ctx context.Context, query *GetOrderByIdQuery) (*dtos.GetOrderByIdResponseDto, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GetOrderByIdHandler.Handle")
	defer span.Finish()

	order, err := q.aggregateStore.Load(ctx, query.OrderId)
	if err != nil {
		return nil, httpErrors.NewNotFoundError(err, fmt.Sprintf("order with id %s not found", order.Id))
	}

	i := order.ShopItems()
	i1 := *i[0]

	s, err := mapper.Map[ordersDto.ShopItemDto](i1)
	fmt.Print(s)

	orderDto, err := mapper.Map[*ordersDto.OrderDto](order)
	if err != nil {
		return nil, err
	}

	return &dtos.GetOrderByIdResponseDto{Order: orderDto}, nil
}
