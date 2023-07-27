package configurations

import (
	"context"

	grpcServer "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/grpc"
	customEcho "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/custom_echo"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/messaging/bus"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orders/internal/orders/configurations/endpoints"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orders/internal/orders/configurations/grpc"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orders/internal/orders/configurations/mappings"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orders/internal/orders/configurations/mediatr"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orders/internal/orders/contracts"
	contracts2 "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orders/internal/shared/contracts"
)

type ordersModuleConfigurator struct {
	*contracts2.InfrastructureConfigurations
	routeBuilder       *customEcho.RouteBuilder
	grpcServiceBuilder *grpcServer.GrpcServiceBuilder
	bus                bus.Bus
	ordersMetrics      *contracts2.OrdersMetrics
}

func NewOrdersModuleConfigurator(infrastructure *contracts2.InfrastructureConfigurations, ordersMetrics *contracts2.OrdersMetrics, bus bus.Bus, routeBuilder *customEcho.RouteBuilder, grpcServiceBuilder *grpcServer.GrpcServiceBuilder) contracts.OrdersModuleConfigurator {
	return &ordersModuleConfigurator{InfrastructureConfigurations: infrastructure, routeBuilder: routeBuilder, grpcServiceBuilder: grpcServiceBuilder, bus: bus, ordersMetrics: ordersMetrics}
}

func (c *ordersModuleConfigurator) ConfigureOrdersModule(ctx context.Context) error {
	//Config Orders Mappings
	err := mappings.ConfigureOrdersMappings()
	if err != nil {
		return err
	}

	//Config Orders Mediators
	err = mediatr.ConfigOrdersMediator(c.InfrastructureConfigurations)
	if err != nil {
		return err
	}

	//Config Orders Grpc
	grpc.ConfigOrdersGrpc(ctx, c.grpcServiceBuilder, c.InfrastructureConfigurations, c.bus, c.ordersMetrics)

	//Config Orders Endpoints
	endpoints.ConfigOrdersEndpoints(ctx, c.routeBuilder, c.InfrastructureConfigurations, c.bus, c.ordersMetrics)

	return nil
}
