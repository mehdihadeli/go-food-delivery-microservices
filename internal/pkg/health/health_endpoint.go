package health

import (
	"net/http"

	customEcho "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/custom_echo"

	"github.com/labstack/echo/v4"
)

type HealthCheckEndpoint struct {
	service    HealthService
	echoServer customEcho.EchoHttpServer
}

func NewHealthCheckEndpoint(
	service HealthService,
	server customEcho.EchoHttpServer,
) *HealthCheckEndpoint {
	return &HealthCheckEndpoint{service: service, echoServer: server}
}

func (s *HealthCheckEndpoint) RegisterEndpoints() {
	s.echoServer.GetEchoInstance().Group("").GET("health", s.CheckHealth)
}

func (s *HealthCheckEndpoint) CheckHealth(c echo.Context) error {
	check := s.service.CheckHealth(c.Request().Context())
	if !check.AllUp() {
		return c.JSON(http.StatusServiceUnavailable, check)
	}
	err := c.JSON(http.StatusOK, check)
	return err
}
