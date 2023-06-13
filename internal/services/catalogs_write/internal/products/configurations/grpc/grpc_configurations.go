package grpc

import (
	googleGrpc "google.golang.org/grpc"

	grpcServer "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/grpc"
	logger2 "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
	productService "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/write_service/internal/products/contracts/proto/service_clients"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/write_service/internal/products/delivery/grpc"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/write_service/internal/shared/contracts"
)

func ConfigProductsGrpc(
	builder *grpcServer.GrpcServiceBuilder,
	logger logger2.Logger,
	metrics *contracts.CatalogsMetrics,
) {
	productGrpcService := grpc.NewProductGrpcService(metrics, logger)
	builder.RegisterRoutes(func(server *googleGrpc.Server) {
		productService.RegisterProductsServiceServer(server, productGrpcService)
	})
}
