package catalogs

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/configurations"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/shared/configurations/infrastructure"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/shared/web"
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

	err = c.migrateCatalogs(infrastructureConfigurations.Gorm)
	if err != nil {
		return err, nil
	}

	c.echo.GET("", func(ec echo.Context) error {
		return ec.String(http.StatusOK, fmt.Sprintf("%s is running...", web.GetMicroserviceName(c.cfg)))
	})

	return nil, infraCleanup
}
