package configurations

import (
	"context"

	productService "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/contracts/proto/service_clients"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/delivery/grpc"
)

func (c *productsModuleConfigurator) configGrpc(ctx context.Context) {
	productGrpcService := grpc.NewProductGrpcService(c.InfrastructureConfiguration)
	productService.RegisterProductsServiceServer(c.GrpcServer.GetCurrentGrpcServer(), productGrpcService)
}
