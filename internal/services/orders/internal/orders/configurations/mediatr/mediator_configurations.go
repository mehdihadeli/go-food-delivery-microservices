package mediatr

import (
	"github.com/mehdihadeli/go-mediatr"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/eventstroredb"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orders/internal/orders/data/repositories"
	createOrderCommandV1 "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orders/internal/orders/features/creating_order/v1/commands"
	createOrderDtosV1 "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orders/internal/orders/features/creating_order/v1/dtos"
	getOrderByIdDtosV1 "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orders/internal/orders/features/getting_order_by_id/v1/dtos"
	getOrderByIdQueryV1 "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orders/internal/orders/features/getting_order_by_id/v1/queries"
	getOrdersDtosV1 "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orders/internal/orders/features/getting_orders/v1/dtos"
	getOrdersQueryV1 "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orders/internal/orders/features/getting_orders/v1/queries"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orders/internal/orders/models/orders/aggregate"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orders/internal/shared/contracts"
)

func ConfigOrdersMediator(infra *contracts.InfrastructureConfigurations) error {
	eventStore := eventstroredb.NewEventStoreDbEventStore(infra.Log, infra.Esdb, infra.EsdbSerializer)
	orderAggregateStore := eventstroredb.NewEventStoreAggregateStore[*aggregate.Order](infra.Log, eventStore, infra.EsdbSerializer)

	mongoOrderReadRepository := repositories.NewMongoOrderReadRepository(infra.Log, infra.Cfg, infra.MongoClient)

	// https://stackoverflow.com/questions/72034479/how-to-implement-generic-interfaces
	err := mediatr.RegisterRequestHandler[*createOrderCommandV1.CreateOrder, *createOrderDtosV1.CreateOrderResponseDto](createOrderCommandV1.NewCreateOrderHandler(infra.Log, infra.Cfg, orderAggregateStore))
	if err != nil {
		return err
	}

	err = mediatr.RegisterRequestHandler[*getOrderByIdQueryV1.GetOrderById, *getOrderByIdDtosV1.GetOrderByIdResponseDto](getOrderByIdQueryV1.NewGetOrderByIdHandler(infra.Log, infra.Cfg, mongoOrderReadRepository))
	if err != nil {
		return err
	}

	err = mediatr.RegisterRequestHandler[*getOrdersQueryV1.GetOrders, *getOrdersDtosV1.GetOrdersResponseDto](getOrdersQueryV1.NewGetOrdersHandler(infra.Log, infra.Cfg, mongoOrderReadRepository))
	if err != nil {
		return err
	}

	return nil
}
