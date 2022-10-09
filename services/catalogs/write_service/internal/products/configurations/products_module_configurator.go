package configurations

import (
	"context"
	grpcServer "github.com/mehdihadeli/store-golang-microservice-sample/pkg/grpc"
	customEcho "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/custom_echo"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/bus"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/configurations/endpoints"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/configurations/grpc"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/configurations/mappings"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/configurations/mediatr"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/contracts"
	repositoriesImp "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/data/repositories"
	contracts2 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/shared/contracts"
)

type productsModuleConfigurator struct {
	contracts2.InfrastructureConfigurations
	routeBuilder       *customEcho.RouteBuilder
	grpcServiceBuilder *grpcServer.GrpcServiceBuilder
	bus                bus.Bus
	catalogsMetrics    contracts2.CatalogsMetrics
}

func NewProductsModuleConfigurator(infrastructure contracts2.InfrastructureConfigurations, catalogsMetrics contracts2.CatalogsMetrics, bus bus.Bus, routeBuilder *customEcho.RouteBuilder, grpcServiceBuilder *grpcServer.GrpcServiceBuilder) contracts.ProductsModuleConfigurator {
	return &productsModuleConfigurator{InfrastructureConfigurations: infrastructure, routeBuilder: routeBuilder, grpcServiceBuilder: grpcServiceBuilder, bus: bus, catalogsMetrics: catalogsMetrics}
}

func (c *productsModuleConfigurator) ConfigureProductsModule(ctx context.Context) error {
	//Config Products Grpc
	grpc.ConfigProductsGrpc(ctx, c.grpcServiceBuilder, c.InfrastructureConfigurations, c.bus, c.catalogsMetrics)

	//Config Products Endpoints
	endpoints.ConfigProductsEndpoints(ctx, c.routeBuilder, c.InfrastructureConfigurations, c.bus, c.catalogsMetrics)

	productRepository := repositoriesImp.NewPostgresProductRepository(c.Log(), c.Cfg(), c.Gorm().DB)

	//Config Products Mappings
	err := mappings.ConfigureProductsMappings()
	if err != nil {
		return err
	}

	//Config Products Mediators
	err = mediatr.ConfigProductsMediator(productRepository, c.InfrastructureConfigurations, c.bus)
	if err != nil {
		return err
	}

	return nil
}
