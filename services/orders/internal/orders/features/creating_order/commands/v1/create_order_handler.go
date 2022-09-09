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

type CreateOrderCommandHandler struct {
	log            logger.Logger
	cfg            *config.Config
	aggregateStore store.AggregateStore[*aggregate.Order]
}

func NewCreateOrderCommandHandler(log logger.Logger, cfg *config.Config, aggregateStore store.AggregateStore[*aggregate.Order]) *CreateOrderCommandHandler {
	return &CreateOrderCommandHandler{log: log, cfg: cfg, aggregateStore: aggregateStore}
}

func (c *CreateOrderCommandHandler) Handle(ctx context.Context, command *CreateOrderCommand) (*dtos.CreateOrderResponseDto, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CreateOrderCommandHandler.Handle")
	span.LogFields(log.String("ProductId", command.OrderID.String()))
	span.LogFields(log.Object("Command", command))
	defer span.Finish()

	shopItems, err := mapper.Map[[]*value_objects.ShopItem](command.ShopItems)
	if err != nil {
		return nil, tracing.TraceWithErr(span, customErrors.NewApplicationErrorWrap(err, "[CreateOrderCommandHandler_Handle.Map] error in the mapping shopItems"))
	}

	order, err := aggregate.NewOrder(command.OrderID, shopItems, command.AccountEmail, command.DeliveryAddress, command.DeliveryTime, command.CreatedAt)

	if err != nil {
		return nil, tracing.TraceWithErr(span, customErrors.NewApplicationErrorWrap(err, "[CreateOrderCommandHandler_Handle.NewOrder] error in creating new order"))
	}

	_, err = c.aggregateStore.Store(order, nil, ctx)
	if err != nil {
		return nil, tracing.TraceWithErr(span, customErrors.NewApplicationErrorWrap(err, "[CreateOrderCommandHandler_Handle.Store] error in storing order aggregate"))
	}

	response := &dtos.CreateOrderResponseDto{OrderID: order.Id()}
	span.LogFields(log.Object("CreateOrderResponseDto", response))

	c.log.Infow(fmt.Sprintf("[CreateOrderCommandHandler.Handle] order with id: {%s} created", command.OrderID), logger.Fields{"ProductId": command.OrderID})

	return response, nil
}
