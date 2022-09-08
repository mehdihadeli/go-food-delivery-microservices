package mediatr

import (
	"github.com/mehdihadeli/go-mediatr"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/es/contracts/store"
	creatingOrderV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/features/creating_order/commands/v1"
	creatingOrderDtos "github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/features/creating_order/dtos"
	gettingOrderByIdDtos "github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/features/getting_order_by_id/dtos"
	gettingOrderByIdV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/features/getting_order_by_id/queries/v1"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/models/orders/aggregate"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/shared/configurations/infrastructure"
)

func ConfigOrdersMediator(aggregateStore store.AggregateStore[*aggregate.Order], infra *infrastructure.InfrastructureConfiguration) error {
	//https://stackoverflow.com/questions/72034479/how-to-implement-generic-interfaces
	err := mediatr.RegisterRequestHandler[*creatingOrderV1.CreateOrderCommand, *creatingOrderDtos.CreateOrderResponseDto](creatingOrderV1.NewCreateOrderHandler(infra.Log, infra.Cfg, aggregateStore))
	if err != nil {
		return err
	}

	err = mediatr.RegisterRequestHandler[*gettingOrderByIdV1.GetOrderByIdQuery, *gettingOrderByIdDtos.GetOrderByIdResponseDto](gettingOrderByIdV1.NewGetOrderByIdHandler(infra.Log, infra.Cfg, aggregateStore))
	if err != nil {
		return err
	}

	return nil
}
