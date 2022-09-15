package mediatr

import (
	"github.com/mehdihadeli/go-mediatr"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/eventstroredb"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/data/repositories"
	creatingOrderV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/features/creating_order/commands/v1"
	creatingOrderDtos "github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/features/creating_order/dtos"
	gettingOrderByIdDtos "github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/features/getting_order_by_id/dtos"
	gettingOrderByIdV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/features/getting_order_by_id/queries/v1"
	gettingOrdersDtos "github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/features/getting_orders/dtos"
	gettingOrdersV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/features/getting_orders/queryies/v1"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/models/orders/aggregate"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/shared/configurations/infrastructure"
)

func ConfigOrdersMediator(infra *infrastructure.InfrastructureConfiguration) error {
	eventStore := eventstroredb.NewEventStoreDbEventStore(infra.Log, infra.Esdb, infra.EsdbSerializer)
	orderAggregateStore := eventstroredb.NewEventStoreAggregateStore[*aggregate.Order](infra.Log, eventStore, infra.EsdbSerializer)

	mongoOrderReadRepository := repositories.NewMongoOrderReadRepository(infra.Log, infra.Cfg, infra.MongoClient)

	//https://stackoverflow.com/questions/72034479/how-to-implement-generic-interfaces
	err := mediatr.RegisterRequestHandler[*creatingOrderV1.CreateOrder, *creatingOrderDtos.CreateOrderResponseDto](creatingOrderV1.NewCreateOrderHandler(infra.Log, infra.Cfg, orderAggregateStore))
	if err != nil {
		return err
	}

	err = mediatr.RegisterRequestHandler[*gettingOrderByIdV1.GetOrderById, *gettingOrderByIdDtos.GetOrderByIdResponseDto](gettingOrderByIdV1.NewGetOrderByIdHandler(infra.Log, infra.Cfg, mongoOrderReadRepository))
	if err != nil {
		return err
	}

	err = mediatr.RegisterRequestHandler[*gettingOrdersV1.GetOrders, *gettingOrdersDtos.GetOrdersResponseDto](gettingOrdersV1.NewGetOrdersHandler(infra.Log, infra.Cfg, mongoOrderReadRepository))
	if err != nil {
		return err
	}

	return nil
}
