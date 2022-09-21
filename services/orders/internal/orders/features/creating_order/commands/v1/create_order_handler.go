package v1

import (
	"context"
	"fmt"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/es/contracts/store"
	customErrors "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/http_errors/custom_errors"
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

type CreateOrderHandler struct {
	log            logger.Logger
	cfg            *config.Config
	aggregateStore store.AggregateStore[*aggregate.Order]
}

func NewCreateOrderHandler(log logger.Logger, cfg *config.Config, aggregateStore store.AggregateStore[*aggregate.Order]) *CreateOrderHandler {
	return &CreateOrderHandler{log: log, cfg: cfg, aggregateStore: aggregateStore}
}

func (c *CreateOrderHandler) Handle(ctx context.Context, command *CreateOrder) (*dtos.CreateOrderResponseDto, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CreateOrderHandler.Handle")
	span.LogFields(log.String("ProductId", command.OrderId.String()))
	span.LogFields(log.Object("Command", command))
	defer span.Finish()

	shopItems, err := mapper.Map[[]*value_objects.ShopItem](command.ShopItems)
	if err != nil {
		return nil, tracing.TraceWithErr(span, customErrors.NewApplicationErrorWrap(err, "[CreateOrderHandler_Handle.Map] error in the mapping shopItems"))
	}

	order, err := aggregate.NewOrder(command.OrderId, shopItems, command.AccountEmail, command.DeliveryAddress, command.DeliveryTime, command.CreatedAt)

	if err != nil {
		return nil, tracing.TraceWithErr(span, customErrors.NewApplicationErrorWrap(err, "[CreateOrderHandler_Handle.NewOrder] error in creating new order"))
	}

	_, err = c.aggregateStore.Store(order, nil, ctx)
	if err != nil {
		return nil, tracing.TraceWithErr(span, customErrors.NewApplicationErrorWrap(err, "[CreateOrderHandler_Handle.Store] error in storing order aggregate"))
	}

	response := &dtos.CreateOrderResponseDto{OrderId: order.Id()}
	span.LogFields(log.Object("CreateOrderResponseDto", response))

	c.log.Infow(fmt.Sprintf("[CreateOrderHandler.Handle] order with id: {%s} created", command.OrderId), logger.Fields{"ProductId": command.OrderId})

	return response, nil
}
