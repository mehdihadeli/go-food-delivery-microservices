package customEcho

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/constants"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/custom_echo/custom_hadnlers"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"strings"
)

type echoHttpServer struct {
	echo   *echo.Echo
	config *EchoHttpConfig
	log    logger.Logger
}

type EchoHttpServer interface {
	RunHttpServer(configEcho func(echoServer *echo.Echo)) error
	GracefulShutdown(ctx context.Context) error
	ApplyVersioningFromHeader()
	GetEchoInstance() *echo.Echo
	SetupDefaultMiddlewares()
	AddMiddlewares(middlewares ...echo.MiddlewareFunc)
	ConfigGroup(groupName string, groupFunc func(group *echo.Group))
}

func NewEchoHttpServer(config *EchoHttpConfig, logger logger.Logger) *echoHttpServer {
	return &echoHttpServer{echo: echo.New(), config: config, log: logger}
}

func (s *echoHttpServer) RunHttpServer(configEcho func(echo *echo.Echo)) error {
	s.echo.Server.ReadTimeout = constants.ReadTimeout
	s.echo.Server.WriteTimeout = constants.WriteTimeout
	s.echo.Server.MaxHeaderBytes = constants.MaxHeaderBytes

	if configEcho != nil {
		configEcho(s.echo)
	}

	//https://echo.labstack.com/guide/http_server/
	return s.echo.Start(s.config.Port)
}

func (s *echoHttpServer) ConfigGroup(groupName string, groupFunc func(group *echo.Group)) {
	groupFunc(s.echo.Group(groupName))
}

func (s *echoHttpServer) AddMiddlewares(middlewares ...echo.MiddlewareFunc) {
	if len(middlewares) > 0 {
		s.echo.Use(middlewares...)
	}
}

func (s *echoHttpServer) GracefulShutdown(ctx context.Context) error {
	err := s.echo.Shutdown(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (s *echoHttpServer) SetupDefaultMiddlewares() {

	s.echo.HideBanner = false
	s.echo.HTTPErrorHandler = customHadnlers.ProblemHandler

	s.echo.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize:         constants.StackSize,
		DisablePrintStack: true,
		DisableStackAll:   true,
	}))
	s.echo.Use(middleware.RequestID())
	s.echo.Use(middleware.Logger())
	s.echo.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: constants.GzipLevel,
		Skipper: func(c echo.Context) bool {
			return strings.Contains(c.Request().URL.Path, "swagger")
		},
	}))

	s.echo.Use(middleware.BodyLimit(constants.BodyLimit))
}

func (s *echoHttpServer) ApplyVersioningFromHeader() {
	s.echo.Pre(apiVersion)
}

func (s *echoHttpServer) GetEchoInstance() *echo.Echo {
	return s.echo
}

// APIVersion Header Based Versioning
func apiVersion(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()
		headers := req.Header

		apiVersion := headers.Get("version")

		req.URL.Path = fmt.Sprintf("/%s%s", apiVersion, req.URL.Path)

		return next(c)
	}
}
