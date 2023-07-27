package configurations

import (
	"context"

	grpcServer "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/grpc"
	customEcho "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/custom_echo"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/messaging/bus"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/write_service/internal/products/configurations/endpoints"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/write_service/internal/products/configurations/grpc"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/write_service/internal/products/configurations/mappings"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/write_service/internal/products/configurations/mediatr"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/write_service/internal/products/contracts"
	repositoriesImp "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/write_service/internal/products/data/repositories"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/write_service/internal/products/data/uow"
	contracts2 "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/write_service/internal/shared/contracts"
)

type productsModuleConfigurator struct {
	*contracts2.InfrastructureConfigurations
	routeBuilder       *customEcho.RouteBuilder
	grpcServiceBuilder *grpcServer.GrpcServiceBuilder
	bus                bus.Bus
	catalogsMetrics    *contracts2.CatalogsMetrics
}

func NewProductsModuleConfigurator(infrastructure *contracts2.InfrastructureConfigurations, catalogsMetrics *contracts2.CatalogsMetrics, bus bus.Bus, routeBuilder *customEcho.RouteBuilder, grpcServiceBuilder *grpcServer.GrpcServiceBuilder) contracts.ProductsModuleConfigurator {
	return &productsModuleConfigurator{InfrastructureConfigurations: infrastructure, routeBuilder: routeBuilder, grpcServiceBuilder: grpcServiceBuilder, bus: bus, catalogsMetrics: catalogsMetrics}
}

func (c *productsModuleConfigurator) ConfigureProductsModule(ctx context.Context) error {
	//cfg Products Grpc
	grpc.ConfigProductsGrpc(ctx, c.grpcServiceBuilder, c.InfrastructureConfigurations, c.bus, c.catalogsMetrics)

	//cfg Products Endpoints
	endpoints.ConfigProductsEndpoints(ctx, c.routeBuilder, c.InfrastructureConfigurations, c.bus, c.catalogsMetrics)

	productRepository := repositoriesImp.NewPostgresProductRepository(c.Log, c.Gorm)
	catalogUnitOfWork := uow.NewCatalogsUnitOfWork(c.Log, c.Gorm)

	//cfg Products Mappings
	err := mappings.ConfigureProductsMappings()
	if err != nil {
		return err
	}

	//cfg Products Mediators
	err = mediatr.ConfigProductsMediator(catalogUnitOfWork, productRepository, c.InfrastructureConfigurations, c.bus)
	if err != nil {
		return err
	}

	return nil
}
