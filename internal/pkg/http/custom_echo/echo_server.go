package customEcho

import (
	"context"
	"fmt"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.opentelemetry.io/otel/metric"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/constants"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/custom_echo/config"
	customHadnlers "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/custom_echo/hadnlers"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/custom_echo/middlewares/log"
	otelMetrics "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/custom_echo/middlewares/otel_metrics"
	otelTracer "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/custom_echo/middlewares/otel_tracer"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
)

type echoHttpServer struct {
	echo         *echo.Echo
	config       *config.EchoHttpOptions
	log          logger.Logger
	meter        metric.Meter
	routeBuilder *RouteBuilder
}

type EchoHttpServer interface {
	RunHttpServer(configEcho ...func(echo *echo.Echo)) error
	GracefulShutdown(ctx context.Context) error
	ApplyVersioningFromHeader()
	GetEchoInstance() *echo.Echo
	Logger() logger.Logger
	Cfg() *config.EchoHttpOptions
	SetupDefaultMiddlewares()
	RouteBuilder() *RouteBuilder
	AddMiddlewares(middlewares ...echo.MiddlewareFunc)
	ConfigGroup(groupName string, groupFunc func(group *echo.Group))
}

func NewEchoHttpServer(
	config *config.EchoHttpOptions,
	logger logger.Logger,
	meter metric.Meter,
) EchoHttpServer {
	e := echo.New()
	e.HideBanner = false

	return &echoHttpServer{
		echo:         e,
		config:       config,
		log:          logger,
		meter:        meter,
		routeBuilder: NewRouteBuilder(e),
	}
}

func (s *echoHttpServer) RunHttpServer(
	configEcho ...func(echo *echo.Echo),
) error {
	s.echo.Server.ReadTimeout = constants.ReadTimeout
	s.echo.Server.WriteTimeout = constants.WriteTimeout
	s.echo.Server.MaxHeaderBytes = constants.MaxHeaderBytes

	if len(configEcho) > 0 {
		ehcoFunc := configEcho[0]
		if ehcoFunc != nil {
			configEcho[0](s.echo)
		}
	}

	// https://echo.labstack.com/guide/http_server/
	return s.echo.Start(s.config.Port)
}

func (s *echoHttpServer) Logger() logger.Logger {
	return s.log
}

func (s *echoHttpServer) Cfg() *config.EchoHttpOptions {
	return s.config
}

func (s *echoHttpServer) RouteBuilder() *RouteBuilder {
	return s.routeBuilder
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
	// set error handler
	s.echo.HTTPErrorHandler = func(err error, c echo.Context) {
		customHadnlers.ProblemHandlerFunc(err, c, s.log)
	}

	// log errors and information
	s.echo.Use(log.EchoLogger(s.log))
	s.echo.Use(otelTracer.Middleware(s.config.Name))
	// Because we use metrics server middleware, if it is not available, our echo will not work.
	if s.meter != nil {
		s.echo.Use(otelMetrics.Middleware(s.meter, s.config.Name))
	}

	s.echo.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogContentLength: true,
		LogLatency:       true,
		LogError:         false,
		LogMethod:        true,
		LogRequestID:     true,
		LogURI:           true,
		LogResponseSize:  true,
		LogURIPath:       true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			s.log.Infow(
				fmt.Sprintf("[Request Middleware] REQUEST: uri: %v, status: %v\n", v.URI, v.Status),
				logger.Fields{"URI": v.URI, "Status": v.Status},
			)
			return nil
		},
	}))
	s.echo.Use(middleware.BodyLimit(constants.BodyLimit))
	s.echo.Use(middleware.RequestID())
	s.echo.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: constants.GzipLevel,
		Skipper: func(c echo.Context) bool {
			return strings.Contains(c.Request().URL.Path, "swagger")
		},
	}))
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
