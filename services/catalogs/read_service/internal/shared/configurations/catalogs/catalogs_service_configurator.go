package catalogs

import (
	"context"
	"emperror.dev/errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/grpc"
	customEcho "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/custom_echo"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/configurations"
	catalogsMetrics "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/shared/configurations/catalogs/metrics"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/shared/configurations/catalogs/rabbitmq"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/shared/contracts"
	"net/http"
)

type catalogsServiceConfigurator struct {
	*contracts.InfrastructureConfigurations
}

func NewCatalogsServiceConfigurator(infra *contracts.InfrastructureConfigurations) contracts.CatalogsServiceConfigurator {
	return &catalogsServiceConfigurator{InfrastructureConfigurations: infra}
}

func (ic *catalogsServiceConfigurator) ConfigureCatalogsService(ctx context.Context) (*contracts.CatalogsServiceConfigurations, error) {
	catalogsServiceConfigurations := &contracts.CatalogsServiceConfigurations{}

	catalogsServiceConfigurations.CatalogsGrpcServer = grpc.NewGrpcServer(ic.Cfg.GRPC, ic.Log, ic.Cfg.ServiceName, ic.Metrics)
	catalogsServiceConfigurations.CatalogsEchoServer = customEcho.NewEchoHttpServer(ic.Cfg.Http, ic.Log, ic.Cfg.ServiceName, ic.Metrics)
	catalogsServiceConfigurations.CatalogsEchoServer.SetupDefaultMiddlewares()

	catalogsServiceConfigurations.CatalogsEchoServer.RouteBuilder().RegisterRoutes(func(e *echo.Echo) {
		e.GET("", func(ec echo.Context) error {
			return ec.String(http.StatusOK, fmt.Sprintf("%s is running...", ic.Cfg.GetMicroserviceNameUpper()))
		})
	})

	// Catalogs Swagger Configs
	ic.configSwagger(catalogsServiceConfigurations.CatalogsEchoServer.RouteBuilder())

	// Catalogs Metrics Configs
	catalogsMetrics, err := catalogsMetrics.ConfigCatalogsMetrics(ic.Cfg, ic.Metrics)
	if err != nil {
		return nil, err
	}
	catalogsServiceConfigurations.CatalogsMetrics = catalogsMetrics

	// Catalogs RabbitMQ Configs
	bus, err := rabbitmq.ConfigCatalogsRabbitMQ(ctx, ic.Cfg.RabbitMQ, ic.InfrastructureConfigurations)
	if err != nil {
		return nil, err
	}
	catalogsServiceConfigurations.CatalogsBus = bus

	// Catalogs Product Module Configs
	pc := configurations.NewProductsModuleConfigurator(ic.InfrastructureConfigurations, catalogsMetrics, bus, catalogsServiceConfigurations.CatalogsEchoServer.RouteBuilder(), catalogsServiceConfigurations.CatalogsGrpcServer.GrpcServiceBuilder())
	err = pc.ConfigureProductsModule(ctx)
	if err != nil {
		return nil, errors.WithMessage(err, "[CatalogsServiceConfigurator_ConfigureCatalogsService.ConfigureProductsModule] error in product module configurator")
	}

	return catalogsServiceConfigurations, nil
}
