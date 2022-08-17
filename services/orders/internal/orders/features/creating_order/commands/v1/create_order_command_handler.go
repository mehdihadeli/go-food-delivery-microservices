package v1

import (
	"context"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/es/contracts/store"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/mapper"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/tracing"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/features/creating_order/dtos"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/models/orders/aggregate"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/models/orders/value_objects"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type CreateOrderCommandHandler struct {
	log            logger.Logger
	cfg            *config.Config
	aggregateStore store.AggregateStore[*aggregate.Order]
}

func NewCreateOrderHandler(log logger.Logger, cfg *config.Config, aggregateStore store.AggregateStore[*aggregate.Order]) *CreateOrderCommandHandler {
	return &CreateOrderCommandHandler{log: log, cfg: cfg, aggregateStore: aggregateStore}
}

func (c *CreateOrderCommandHandler) Handle(ctx context.Context, command *CreateOrderCommand) (*dtos.CreateOrderResponseDto, error) {

	span, ctx := opentracing.StartSpanFromContext(ctx, "CreateOrderCommandHandler.Handle")
	defer span.Finish()
	span.LogFields(log.Object("command", command))

	shopItems, err := mapper.Map[[]*value_objects.ShopItem](command.ShopItems)
	if err != nil {
		return nil, err
	}

	order, err := aggregate.NewOrder(command.OrderID, shopItems, command.AccountEmail, command.DeliveryAddress, command.DeliveryTime, command.CreatedAt)

	if err != nil {
		return nil, tracing.TraceWithErr(span, err)
	}
	_, err = c.aggregateStore.Store(order, nil, ctx)
	if err != nil {
		return nil, err
	}

	return &dtos.CreateOrderResponseDto{OrderID: order.Id()}, nil
}
