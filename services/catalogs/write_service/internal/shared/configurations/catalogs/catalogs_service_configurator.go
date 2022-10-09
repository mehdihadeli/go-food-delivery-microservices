package catalogs

import (
	"context"
	"emperror.dev/errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/grpc"
	customEcho "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/custom_echo"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/configurations"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/shared/configurations/catalogs/metrics"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/shared/configurations/catalogs/rabbitmq"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/shared/contracts"
	"net/http"
)

type catalogsServiceConfigurator struct {
	contracts.InfrastructureConfigurations
}

func NewCatalogsServiceConfigurator(infrastructureConfiguration contracts.InfrastructureConfigurations) contracts.CatalogsServiceConfigurator {
	return &catalogsServiceConfigurator{InfrastructureConfigurations: infrastructureConfiguration}
}

func (c *catalogsServiceConfigurator) ConfigureCatalogsService(ctx context.Context) (contracts.CatalogServiceConfigurations, error) {
	catalogsServiceConfigurations := &catalogsServiceConfigurations{}

	catalogsServiceConfigurations.catalogsGrpcServer = grpc.NewGrpcServer(c.Cfg().GRPC, c.Log(), c.Cfg().ServiceName, c.Metrics())
	catalogsServiceConfigurations.catalogsEchoServer = customEcho.NewEchoHttpServer(c.Cfg().Http, c.Log(), c.Cfg().ServiceName, c.Metrics())

	catalogsServiceConfigurations.catalogsEchoServer.RouteBuilder().RegisterRoutes(func(e *echo.Echo) {
		e.GET("", func(ec echo.Context) error {
			return ec.String(http.StatusOK, fmt.Sprintf("%s is running...", c.Cfg().GetMicroserviceNameUpper()))
		})
	})

	// Catalogs Swagger Configs
	c.configSwagger(catalogsServiceConfigurations.catalogsEchoServer.RouteBuilder())

	// Catalogs Metrics Configs
	metrics, err := metrics.ConfigCatalogsMetrics(c.Cfg(), c.Metrics())
	if err != nil {
		return nil, err
	}
	catalogsServiceConfigurations.catalogsMetrics = metrics

	// Catalogs RabbitMQ Configs
	bus, err := rabbitmq.ConfigCatalogsRabbitMQ(ctx, c.Cfg().RabbitMQ, c.InfrastructureConfigurations)
	if err != nil {
		return nil, err
	}
	catalogsServiceConfigurations.catalogsBus = bus

	// Catalogs Product Module Configs
	pc := configurations.NewProductsModuleConfigurator(c.InfrastructureConfigurations, metrics, bus, catalogsServiceConfigurations.catalogsEchoServer.RouteBuilder(), catalogsServiceConfigurations.catalogsGrpcServer.GrpcServiceBuilder())
	err = pc.ConfigureProductsModule(ctx)
	if err != nil {
		return nil, errors.WithMessage(err, "[CatalogsServiceConfigurator_ConfigureCatalogsService.ConfigureProductsModule] error in product module configurator")
	}

	err = c.migrateCatalogs(c.Gorm())
	if err != nil {
		return nil, errors.WithMessage(err, "[CatalogsServiceConfigurator_ConfigureCatalogsService.migrateCatalogs] error in migrateCatalogs")
	}

	return catalogsServiceConfigurations, nil
}
