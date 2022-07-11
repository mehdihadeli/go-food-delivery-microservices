package configurations

import (
	"context"

	product_service "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/contracts/grpc/service_clients"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/delivery/grpc"
)

func (c *productsModuleConfigurator) configGrpc(ctx context.Context) {
	productGrpcService := grpc.NewProductGrpcService(c.InfrastructureConfiguration)
	product_service.RegisterProductsServiceServer(c.GrpcServer, productGrpcService)
}
