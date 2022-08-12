package configurations

import (
	"context"
	esSerializer "github.com/mehdihadeli/store-golang-microservice-sample/pkg/es/serializer"
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

	err := mappings.ConfigureMappings()
	if err != nil {
		return err
	}

	serializer := eventstroredb.NewEsdbSerializer(esSerializer.NewJsonMetadataSerializer(), esSerializer.NewJsonEventSerializer())
	eventStore := eventstroredb.NewEventStoreDbEventStore(c.Log, c.Esdb, serializer)
	aggregateStore := eventstroredb.NewEventStoreAggregateStore[*aggregate.Order](c.Log, eventStore, serializer)

	err = c.configOrdersMediator(aggregateStore)
	if err != nil {
		return err
	}

	if c.Cfg.DeliveryType == "grpc" {
		c.configGrpc(ctx)
	} else {
		c.configEndpoints(ctx)
	}

	return nil
}
