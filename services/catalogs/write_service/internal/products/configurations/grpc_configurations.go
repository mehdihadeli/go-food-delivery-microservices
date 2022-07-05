package configurations

import (
	product_service "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/contracts/grpc/service_clients"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/delivery/grpc"
)

func (pm *ProductModule) configGrpc() {
	productGrpcService := grpc.NewProductGrpcService(pm.Infrastructure, pm.Mediator, pm.ProductRepository)
	product_service.RegisterProductsServiceServer(pm.Infrastructure.GrpcServer, productGrpcService)
}
