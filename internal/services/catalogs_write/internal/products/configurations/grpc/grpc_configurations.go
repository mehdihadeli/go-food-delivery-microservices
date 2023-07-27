package grpc

import (
	"context"

	googleGrpc "google.golang.org/grpc"

	grpcServer "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/grpc"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/messaging/bus"

	productService "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/write_service/internal/products/contracts/proto/service_clients"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/write_service/internal/products/delivery/grpc"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/write_service/internal/shared/contracts"
)

func ConfigProductsGrpc(ctx context.Context, builder *grpcServer.GrpcServiceBuilder, infra *contracts.InfrastructureConfigurations, bus bus.Bus, metrics *contracts.CatalogsMetrics) {
	productGrpcService := grpc.NewProductGrpcService(infra, metrics, bus)
	builder.RegisterRoutes(func(server *googleGrpc.Server) {
		productService.RegisterProductsServiceServer(server, productGrpcService)
	})
}
