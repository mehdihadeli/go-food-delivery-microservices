package server

import (
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/config"
	"google.golang.org/grpc"
	"net/http"
)

type Server struct {
	Log          logger.Logger
	Cfg          *config.Config
	Echo         *echo.Echo
	DoneCh       chan struct{}
	GrpcServer   *grpc.Server
	HealthServer *http.Server
}

func NewServer(log logger.Logger, cfg *config.Config) *Server {

	return &Server{Log: log, Cfg: cfg, Echo: echo.New(), GrpcServer: NewGrpcServer(), HealthServer: NewHealthCheckServer(cfg)}
}
