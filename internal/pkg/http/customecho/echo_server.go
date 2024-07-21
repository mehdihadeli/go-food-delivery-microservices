package customEcho

import (
	"context"
	"fmt"
	"strings"

	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/constants"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/http/customecho/config"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/http/customecho/contracts"
	hadnlers "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/http/customecho/hadnlers"
	ipratelimit "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/http/customecho/middlewares/ip_ratelimit"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/http/customecho/middlewares/log"
	otelMetrics "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/http/customecho/middlewares/otel_metrics"
	oteltracing "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/http/customecho/middlewares/otel_tracing"
	problemdetail "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/http/customecho/middlewares/problem_detail"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/logger"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.opentelemetry.io/otel/metric"
)

type echoHttpServer struct {
	echo         *echo.Echo
	config       *config.EchoHttpOptions
	log          logger.Logger
	meter        metric.Meter
	routeBuilder *contracts.RouteBuilder
}

func NewEchoHttpServer(
	config *config.EchoHttpOptions,
	logger logger.Logger,
	meter metric.Meter,
) contracts.EchoHttpServer {
	e := echo.New()
	e.HideBanner = true

	return &echoHttpServer{
		echo:         e,
		config:       config,
		log:          logger,
		meter:        meter,
		routeBuilder: contracts.NewRouteBuilder(e),
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

func (s *echoHttpServer) RouteBuilder() *contracts.RouteBuilder {
	return s.routeBuilder
}

func (s *echoHttpServer) ConfigGroup(
	groupName string,
	groupFunc func(group *echo.Group),
) {
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
	skipper := func(c echo.Context) bool {
		return strings.Contains(c.Request().URL.Path, "swagger") ||
			strings.Contains(c.Request().URL.Path, "metrics") ||
			strings.Contains(c.Request().URL.Path, "health") ||
			strings.Contains(c.Request().URL.Path, "favicon.ico")
	}

	// set error handler
	s.echo.HTTPErrorHandler = func(err error, c echo.Context) {
		// bypass skip endpoints and its error
		if skipper(c) {
			return
		}

		hadnlers.ProblemDetailErrorHandlerFunc(err, c, s.log)
	}

	// log errors and information
	s.echo.Use(
		log.EchoLogger(
			s.log,
			log.WithSkipper(skipper),
		),
	)
	s.echo.Use(
		oteltracing.HttpTrace(
			oteltracing.WithSkipper(skipper),
			oteltracing.WithServiceName(s.config.Name),
		),
	)
	s.echo.Use(
		otelMetrics.HTTPMetrics(
			otelMetrics.WithServiceName(s.config.Name),
			otelMetrics.WithSkipper(skipper)),
	)
	s.echo.Use(middleware.BodyLimit(constants.BodyLimit))
	s.echo.Use(ipratelimit.IPRateLimit())
	s.echo.Use(middleware.RequestID())
	s.echo.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level:   constants.GzipLevel,
		Skipper: skipper,
	}))
	// should be last middleware
	s.echo.Use(problemdetail.ProblemDetail(problemdetail.WithSkipper(skipper)))
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
