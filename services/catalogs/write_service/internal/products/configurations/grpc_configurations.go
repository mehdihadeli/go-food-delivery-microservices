package configurations

import (
	"context"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/mediatr"
	product_service "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/contracts/grpc/service_clients"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/delivery/grpc"
)

func (c *productsModuleConfigurator) configGrpc(ctx context.Context, mediator *mediatr.Mediator) {
	productGrpcService := grpc.NewProductGrpcService(c.InfrastructureConfiguration, mediator)
	product_service.RegisterProductsServiceServer(c.GrpcServer, productGrpcService)
}
