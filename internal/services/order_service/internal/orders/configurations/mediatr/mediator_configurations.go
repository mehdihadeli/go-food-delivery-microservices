package mediatr

import (
	"github.com/mehdihadeli/go-mediatr"

	repositories2 "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orderservice/internal/orders/contracts/repositories"
	createOrderCommandV1 "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orderservice/internal/orders/features/creating_order/v1/commands"
	createOrderDtosV1 "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orderservice/internal/orders/features/creating_order/v1/dtos"
	getOrderByIdDtosV1 "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orderservice/internal/orders/features/getting_order_by_id/v1/dtos"
	getOrderByIdQueryV1 "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orderservice/internal/orders/features/getting_order_by_id/v1/queries"
	getOrdersDtosV1 "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orderservice/internal/orders/features/getting_orders/v1/dtos"
	getOrdersQueryV1 "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orderservice/internal/orders/features/getting_orders/v1/queries"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orderservice/internal/orders/models/orders/aggregate"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/es/contracts/store"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/otel/tracing"
)

func ConfigOrdersMediator(
	logger logger.Logger,
	mongoOrderReadRepository repositories2.OrderMongoRepository,
	orderAggregateStore store.AggregateStore[*aggregate.Order],
	tracer tracing.AppTracer,
) error {
	// https://stackoverflow.com/questions/72034479/how-to-implement-generic-interfaces
	err := mediatr.RegisterRequestHandler[*createOrderCommandV1.CreateOrder, *createOrderDtosV1.CreateOrderResponseDto](
		createOrderCommandV1.NewCreateOrderHandler(logger, orderAggregateStore, tracer),
	)
	if err != nil {
		return err
	}

	err = mediatr.RegisterRequestHandler[*getOrderByIdQueryV1.GetOrderById, *getOrderByIdDtosV1.GetOrderByIdResponseDto](
		getOrderByIdQueryV1.NewGetOrderByIdHandler(logger, mongoOrderReadRepository, tracer),
	)
	if err != nil {
		return err
	}

	err = mediatr.RegisterRequestHandler[*getOrdersQueryV1.GetOrders, *getOrdersDtosV1.GetOrdersResponseDto](
		getOrdersQueryV1.NewGetOrdersHandler(logger, mongoOrderReadRepository, tracer),
	)
	if err != nil {
		return err
	}

	return nil
}
