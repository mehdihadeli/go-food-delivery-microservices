package order_module

import (
	"context"
	orders_service "github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/contracts/proto/service_clients"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/delivery/grpc"
)

func (c *ordersModuleConfigurator) configGrpc(ctx context.Context) {
	orderGrpcService := grpc.NewOrderGrpcService(c.InfrastructureConfiguration)
	orders_service.RegisterOrdersServiceServer(c.grpcServer.GetCurrentGrpcServer(), orderGrpcService)
}
