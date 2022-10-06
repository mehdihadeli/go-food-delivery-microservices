package order_module

import (
	"context"
	grpcServer "github.com/mehdihadeli/store-golang-microservice-sample/pkg/grpc"
	customEcho "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/custom_echo"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/configurations/mappings"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/configurations/mediatr"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/configurations/projections"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/contracts"
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

	err = mediatr.ConfigOrdersMediator(c.InfrastructureConfiguration)
	if err != nil {
		return err
	}

	if c.Cfg.DeliveryType == "grpc" {
		c.configGrpc(ctx)
	} else {
		c.configEndpoints(ctx)
	}

	projections.ConfigOrderProjections(c.InfrastructureConfiguration)

	return nil
}
