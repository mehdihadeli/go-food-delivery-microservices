package grpc

import (
	"context"
	grpcServer "github.com/mehdihadeli/store-golang-microservice-sample/pkg/grpc"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/bus"
	productService "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/contracts/proto/service_clients"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/delivery/grpc"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/shared/contracts"
	googleGrpc "google.golang.org/grpc"
)

func ConfigProductsGrpc(ctx context.Context, builder *grpcServer.GrpcServiceBuilder, infra *contracts.InfrastructureConfigurations, bus bus.Bus, metrics *contracts.CatalogsMetrics) {
	productGrpcService := grpc.NewProductGrpcService(infra, metrics, bus)
	builder.RegisterRoutes(func(server *googleGrpc.Server) {
		productService.RegisterProductsServiceServer(server, productGrpcService)
	})
}
