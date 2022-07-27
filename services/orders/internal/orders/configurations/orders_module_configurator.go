package configurations

import (
	"context"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/eventstroredb"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/mappings"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/models/orders/aggregate"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/shared/configurations/infrastructure"
)

type OrdersModuleConfigurator interface {
	ConfigureProductsModule() error
}

type ordersModuleConfigurator struct {
	*infrastructure.InfrastructureConfiguration
}

func NewOrdersModuleConfigurator(infrastructure *infrastructure.InfrastructureConfiguration) *ordersModuleConfigurator {
	return &ordersModuleConfigurator{InfrastructureConfiguration: infrastructure}
}

func (c *ordersModuleConfigurator) ConfigureOrdersModule(ctx context.Context) error {

	v1 := c.Echo.Group("/api/v1")
	group := v1.Group("/" + c.Cfg.Http.OrdersPath)

	err := mappings.ConfigureMappings()
	if err != nil {
		return err
	}

	aggregateStore := eventstroredb.NewEventStoreAggregateStore[*aggregate.Order](c.Log, c.Esdb)
	err = c.configOrdersMediator(aggregateStore)
	if err != nil {
		return err
	}

	c.configEndpoints(ctx, group)

	if c.Cfg.DeliveryType == "grpc" {
		c.configGrpc(ctx)
	}

	return nil
}
