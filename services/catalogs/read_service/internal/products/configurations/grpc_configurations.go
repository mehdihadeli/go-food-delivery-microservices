package configurations

import (
	"fmt"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/mediatr"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/shared/configurations"
)

type ProductGrpcConfigurator struct {
	*ProductModuleConfigurations
}

type ProductGrpcConfigurations struct {
	*configurations.Infrastructure
	*mediatr.Mediator
}

func (pc *ProductGrpcConfigurator) configGrpc() {
	grpcConfigurations := &ProductGrpcConfigurations{Infrastructure: pc.Infrastructure, Mediator: pc.Mediator}
	fmt.Print(grpcConfigurations)
	//productGrpcService := grpc.NewProductGrpcService(pm.Infrastructure, pm.Mediator, pm.ProductRepository)
	//product_service.RegisterProductsServiceServer(pm.Infrastructure.GrpcServer, productGrpcService)
}
