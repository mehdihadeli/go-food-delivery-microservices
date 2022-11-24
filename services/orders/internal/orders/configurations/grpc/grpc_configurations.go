package grpc

import (
	"context"

	googleGrpc "google.golang.org/grpc"

	grpcServer "github.com/mehdihadeli/store-golang-microservice-sample/pkg/grpc"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/bus"
	ordersService "github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/contracts/proto/service_clients"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/delivery/grpc"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/shared/contracts"
)

func ConfigOrdersGrpc(ctx context.Context, builder *grpcServer.GrpcServiceBuilder, infra contracts.InfrastructureConfigurations, bus bus.Bus, metrics contracts.OrdersMetrics) {
	orderGrpcService := grpc.NewOrderGrpcService(infra, metrics, bus)
	builder.RegisterRoutes(func(server *googleGrpc.Server) {
		ordersService.RegisterOrdersServiceServer(server, orderGrpcService)
	})
}
