package server

import (
	"context"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/grpc"
	customEcho "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/custom_echo"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger/defaultLogger"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/shared/configurations/catalogs"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/shared/configurations/infrastructure"
)

type TestServer struct {
	Log        logger.Logger
	Cfg        *config.Config
	HttpServer customEcho.EchoHttpServer
	GrpcServer grpc.GrpcServer
}

func NewTestServer() *TestServer {
	cfg, err := config.InitConfig("test")

	c := infrastructure.NewInfrastructureConfigurator(defaultLogger.Logger, cfg)
	infrastructures, err, cleanup := c.ConfigInfrastructures(context.Background())
	if err != nil {
		return nil
	}

	defer cleanup()

	httpServer := customEcho.NewEchoHttpServer(cfg.Http, defaultLogger.Logger)
	grpcServer := grpc.NewGrpcServer(cfg.GRPC, defaultLogger.Logger)

	catalog := catalogs.NewCatalogsServiceConfigurator(infrastructures, httpServer, grpcServer)
	err = catalog.ConfigureCatalogsService(context.Background())
	if err != nil {
		return nil
	}

	return &TestServer{
		GrpcServer: grpcServer,
		HttpServer: httpServer,
		Log:        defaultLogger.Logger,
		Cfg:        cfg,
	}
}
