package configurations

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/es/contracts/store"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/mediatr"
	creatingOrderV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/features/creating_order/commands/v1"
	creatingOrderDtos "github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/features/creating_order/dtos"
	gettingOrderByIdDtos "github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/features/getting_order_by_id/dtos"
	gettingOrderByIdV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/features/getting_order_by_id/queries/v1"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/models/orders/aggregate"
)

func (c *ordersModuleConfigurator) configOrdersMediator(aggregateStore store.AggregateStore[*aggregate.Order]) error {

	//https://stackoverflow.com/questions/72034479/how-to-implement-generic-interfaces
	err := mediatr.RegisterHandler[*creatingOrderV1.CreateOrderCommand, *creatingOrderDtos.CreateOrderResponseDto](creatingOrderV1.NewCreateOrderHandler(c.Log, c.Cfg, aggregateStore))
	if err != nil {
		return err
	}

	err = mediatr.RegisterHandler[*gettingOrderByIdV1.GetOrderByIdQuery, *gettingOrderByIdDtos.GetOrderByIdResponseDto](gettingOrderByIdV1.NewGetOrderByIdHandler(c.Log, c.Cfg, aggregateStore))
	if err != nil {
		return err
	}

	return nil
}
