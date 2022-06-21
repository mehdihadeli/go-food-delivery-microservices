package configurations

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/docs"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/shared/constants"
	echoSwagger "github.com/swaggo/echo-swagger"
	"strings"
)

func (s *Server) runHttpServer() error {
	s.configSwagger()

	s.Echo.Server.ReadTimeout = constants.ReadTimeout
	s.Echo.Server.WriteTimeout = constants.WriteTimeout
	s.Echo.Server.MaxHeaderBytes = constants.MaxHeaderBytes

	return s.Echo.Start(s.Cfg.Http.Port)
}

func (s *Server) configSwagger() {
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Title = "Catalogs Service Api"
	docs.SwaggerInfo.Description = "Catalogs Service Api."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = "/api/v1"
	s.Echo.GET("/swagger/*", echoSwagger.WrapHandler)

	s.Echo.Use(s.middlewareManager.RequestLoggerMiddleware)
	s.Echo.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize:         constants.StackSize,
		DisablePrintStack: true,
		DisableStackAll:   true,
	}))
	s.Echo.Use(middleware.RequestID())
	s.Echo.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: constants.GzipLevel,
		Skipper: func(c echo.Context) bool {
			return strings.Contains(c.Request().URL.Path, "swagger")
		},
	}))
	s.Echo.Use(middleware.BodyLimit(constants.BodyLimit))
}

func (s *Server) applyVersioningFromHeader() {
	s.Echo.Pre(apiVersion)
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
