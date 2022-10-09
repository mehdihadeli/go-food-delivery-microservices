package grpc

import (
	"context"
	grpcServer "github.com/mehdihadeli/store-golang-microservice-sample/pkg/grpc"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/bus"
	orders_service "github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/contracts/proto/service_clients"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/delivery/grpc"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/shared/contracts"
	googleGrpc "google.golang.org/grpc"
)

func ConfigOrdersGrpc(ctx context.Context, builder *grpcServer.GrpcServiceBuilder, infra contracts.InfrastructureConfigurations, bus bus.Bus, metrics contracts.OrdersMetrics) {
	orderGrpcService := grpc.NewOrderGrpcService(infra, metrics, bus)
	builder.RegisterRoutes(func(server *googleGrpc.Server) {
		orders_service.RegisterOrdersServiceServer(server, orderGrpcService)
	})
}
