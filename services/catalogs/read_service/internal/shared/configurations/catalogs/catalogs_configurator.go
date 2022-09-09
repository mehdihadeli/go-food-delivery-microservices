package catalogs

import (
	"context"
	"emperror.dev/errors"
	"fmt"
	"github.com/labstack/echo/v4"
	grpcServer "github.com/mehdihadeli/store-golang-microservice-sample/pkg/grpc"
	customEcho "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/custom_echo"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/configurations/product_module"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/shared/configurations/infrastructure"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/shared/web"
	"net/http"
)

type CatalogsServiceConfigurator interface {
	ConfigureProductsModule() error
}

type catalogsServiceConfigurator struct {
	*infrastructure.InfrastructureConfigurations
	echoServer customEcho.EchoHttpServer
	grpcServer grpcServer.GrpcServer
}

func NewCatalogsServiceConfigurator(infra *infrastructure.InfrastructureConfigurations, echoServer customEcho.EchoHttpServer, grpcServer grpcServer.GrpcServer) *catalogsServiceConfigurator {
	return &catalogsServiceConfigurator{echoServer: echoServer, grpcServer: grpcServer, InfrastructureConfigurations: infra}
}

func (c *catalogsServiceConfigurator) ConfigureCatalogsService(ctx context.Context) error {
	pc := product_module.NewProductsModuleConfigurator(c.InfrastructureConfigurations, c.echoServer, c.grpcServer)
	err := pc.ConfigureProductsModule(ctx)
	if err != nil {
		return errors.WithMessage(err, "[CatalogsServiceConfigurator_ConfigureCatalogsService.ConfigureProductsModule] error in product module configurator")
	}

	c.migrationCatalogsMongo(ctx, c.MongoClient)

	c.echoServer.GetEchoInstance().GET("", func(ec echo.Context) error {
		return ec.String(http.StatusOK, fmt.Sprintf("%s is running...", web.GetMicroserviceName(c.Cfg)))
	})

	c.configSwagger()
	c.configMiddlewares(c.Metrics)

	return nil
}
