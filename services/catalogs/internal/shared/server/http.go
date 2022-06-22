package server

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/docs"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/shared/constants"
	echoSwagger "github.com/swaggo/echo-swagger"
	"time"
)

func (s *Server) RunHttpServer() error {
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
}

func (s *Server) ApplyVersioningFromHeader() {
	s.Echo.Pre(apiVersion)
}

func (s *Server) WaitShootDown(duration time.Duration) {
	go func() {
		time.Sleep(duration)
		s.DoneCh <- struct{}{}
	}()
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
