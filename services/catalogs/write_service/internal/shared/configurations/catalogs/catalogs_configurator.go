package catalogs

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	grpcServer "github.com/mehdihadeli/store-golang-microservice-sample/pkg/grpc"
	customEcho "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/custom_echo"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/configurations/product_mdoule"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/shared/configurations/infrastructure"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/shared/web"
	"github.com/pkg/errors"
	"net/http"
)

type CatalogsServiceConfigurator interface {
	ConfigureProductsModule() error
}

type catalogsServiceConfigurator struct {
	*infrastructure.InfrastructureConfiguration
	echoServer customEcho.EchoHttpServer
	grpcServer grpcServer.GrpcServer
}

func NewCatalogsServiceConfigurator(infrastructureConfiguration *infrastructure.InfrastructureConfiguration, echoServer customEcho.EchoHttpServer, grpcServer grpcServer.GrpcServer) *catalogsServiceConfigurator {
	return &catalogsServiceConfigurator{echoServer: echoServer, grpcServer: grpcServer, InfrastructureConfiguration: infrastructureConfiguration}
}

func (c *catalogsServiceConfigurator) ConfigureCatalogsService(ctx context.Context) error {
	pc := product_mdoule.NewProductsModuleConfigurator(c.InfrastructureConfiguration, c.echoServer, c.grpcServer)
	err := pc.ConfigureProductsModule(ctx)
	if err != nil {
		return errors.WithMessage(err, "[CatalogsServiceConfigurator_ConfigureCatalogsService.ConfigureProductsModule] error in product module configurator")
	}

	err = c.migrateCatalogs(c.Gorm)
	if err != nil {
		return errors.WithMessage(err, "[CatalogsServiceConfigurator_ConfigureCatalogsService.migrateCatalogs] error in migrateCatalogs")
	}

	c.configSwagger()
	c.configMiddlewares(c.Metrics)

	c.echoServer.GetEchoInstance().GET("", func(ec echo.Context) error {
		return ec.String(http.StatusOK, fmt.Sprintf("%s is running...", web.GetMicroserviceName(c.Cfg)))
	})

	return nil
}
