package catalogs

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/configurations"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/shared/configurations/infrastructure"
	"google.golang.org/grpc"
	"net/http"
)

type CatalogsServiceConfigurator interface {
	ConfigureProductsModule() error
}

type catalogsServiceConfigurator struct {
	log        logger.Logger
	cfg        *config.Config
	echo       *echo.Echo
	grpcServer *grpc.Server
}

func NewCatalogsServiceConfigurator(log logger.Logger, cfg *config.Config, echo *echo.Echo, grpcServer *grpc.Server) *catalogsServiceConfigurator {
	return &catalogsServiceConfigurator{cfg: cfg, echo: echo, grpcServer: grpcServer, log: log}
}

func (c *catalogsServiceConfigurator) ConfigureCatalogsService(ctx context.Context) (error, func()) {

	ic := infrastructure.NewInfrastructureConfigurator(c.log, c.cfg, c.echo, c.grpcServer)
	infrastructureConfigurations, err, infraCleanup := ic.ConfigInfrastructures(ctx)
	if err != nil {
		return err, nil
	}

	pc := configurations.NewProductsModuleConfigurator(infrastructureConfigurations)
	err = pc.ConfigureProductsModule(ctx)
	if err != nil {
		return err, nil
	}
	c.migrationCatalogsMongo(ctx, infrastructureConfigurations.MongoClient)

	c.echo.GET("", func(c echo.Context) error {
		return c.String(http.StatusOK, "Catalogs Read-Service is running...")
	})

	return nil, infraCleanup
}
