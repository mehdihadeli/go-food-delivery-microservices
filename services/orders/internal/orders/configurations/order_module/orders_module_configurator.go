package order_module

import (
	"context"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/eventstroredb"
	grpcServer "github.com/mehdihadeli/store-golang-microservice-sample/pkg/grpc"
	customEcho "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/custom_echo"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/configurations/mappings"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/configurations/mediatr"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/contracts"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/models/orders/aggregate"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/shared/configurations/infrastructure"
)

type ordersModuleConfigurator struct {
	*infrastructure.InfrastructureConfiguration
	echoServer customEcho.EchoHttpServer
	grpcServer grpcServer.GrpcServer
}

func NewOrdersModuleConfigurator(infrastructure *infrastructure.InfrastructureConfiguration, echoServer customEcho.EchoHttpServer, grpcServer grpcServer.GrpcServer) contracts.OrdersModuleConfigurator {
	return &ordersModuleConfigurator{InfrastructureConfiguration: infrastructure, echoServer: echoServer, grpcServer: grpcServer}
}

func (c *ordersModuleConfigurator) ConfigureOrdersModule(ctx context.Context) error {
	err := mappings.ConfigureMappings()
	if err != nil {
		return err
	}

	eventStore := eventstroredb.NewEventStoreDbEventStore(c.Log, c.Esdb, c.EsdbSerializer)
	aggregateStore := eventstroredb.NewEventStoreAggregateStore[*aggregate.Order](c.Log, eventStore, c.EsdbSerializer)

	err = mediatr.ConfigOrdersMediator(aggregateStore, c.InfrastructureConfiguration)
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
