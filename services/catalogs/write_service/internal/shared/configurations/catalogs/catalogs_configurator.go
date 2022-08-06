package catalogs

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	grpcServer "github.com/mehdihadeli/store-golang-microservice-sample/pkg/grpc"
	customEcho "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/custom_echo"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/configurations"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/shared/configurations/infrastructure"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/shared/web"
	"net/http"
)

type CatalogsServiceConfigurator interface {
	ConfigureProductsModule() error
}

type catalogsServiceConfigurator struct {
	log        logger.Logger
	cfg        *config.Config
	echoServer customEcho.EchoHttpServer
	grpcServer grpcServer.GrpcServer
}

func NewCatalogsServiceConfigurator(log logger.Logger, cfg *config.Config, echoServer customEcho.EchoHttpServer, grpcServer grpcServer.GrpcServer) *catalogsServiceConfigurator {
	return &catalogsServiceConfigurator{cfg: cfg, echoServer: echoServer, grpcServer: grpcServer, log: log}
}

func (c *catalogsServiceConfigurator) ConfigureCatalogsService(ctx context.Context) (error, func()) {

	ic := infrastructure.NewInfrastructureConfigurator(c.log, c.cfg, c.echoServer, c.grpcServer)
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

	c.echoServer.GetEchoInstance().GET("", func(ec echo.Context) error {
		return ec.String(http.StatusOK, fmt.Sprintf("%s is running...", web.GetMicroserviceName(c.cfg)))
	})

	return nil, infraCleanup
}
