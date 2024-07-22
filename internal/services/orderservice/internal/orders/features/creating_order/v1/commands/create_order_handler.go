package createOrderCommandV1

import (
	"context"
	"fmt"

	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/es/contracts/store"
	customErrors "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/http/httperrors/customerrors"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/logger"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/mapper"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/otel/tracing"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/services/orderservice/internal/orders/features/creating_order/v1/dtos"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/services/orderservice/internal/orders/models/orders/aggregate"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/services/orderservice/internal/orders/models/orders/value_objects"
)

type CreateOrderHandler struct {
	log logger.Logger
	// goland can't detect this generic type, but it is ok in vscode
	aggregateStore store.AggregateStore[*aggregate.Order]
	tracer         tracing.AppTracer
}

func NewCreateOrderHandler(
	log logger.Logger,
	aggregateStore store.AggregateStore[*aggregate.Order],
	tracer tracing.AppTracer,
) *CreateOrderHandler {
	return &CreateOrderHandler{log: log, aggregateStore: aggregateStore, tracer: tracer}
}

func (c *CreateOrderHandler) Handle(
	ctx context.Context,
	command *CreateOrder,
) (*dtos.CreateOrderResponseDto, error) {
	shopItems, err := mapper.Map[[]*value_objects.ShopItem](command.ShopItems)
	if err != nil {
		return nil,
			customErrors.NewApplicationErrorWrap(
				err,
				"[CreateOrderHandler_Handle.Map] error in the mapping shopItems",
			)
	}

	order, err := aggregate.NewOrder(
		command.OrderId,
		shopItems,
		command.AccountEmail,
		command.DeliveryAddress,
		command.DeliveryTime,
		command.CreatedAt,
	)
	if err != nil {
		return nil, customErrors.NewApplicationErrorWrap(
			err,
			"[CreateOrderHandler_Handle.NewOrder] error in creating new order",
		)
	}

	_, err = c.aggregateStore.Store(order, nil, ctx)
	if err != nil {
		return nil, customErrors.NewApplicationErrorWrap(
			err,
			"[CreateOrderHandler_Handle.Store] error in storing order aggregate",
		)
	}

	response := &dtos.CreateOrderResponseDto{OrderId: order.Id()}

	c.log.Infow(
		fmt.Sprintf("[CreateOrderHandler.Handle] order with id: {%s} created", command.OrderId),
		logger.Fields{"ProductId": command.OrderId},
	)

	return response, nil
}
