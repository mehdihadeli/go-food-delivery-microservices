package params

import (
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
	"go.uber.org/fx"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orderservice/internal/shared/contracts"
)

type OrderRouteParams struct {
	fx.In

	OrdersMetrics *contracts.OrdersMetrics
	Logger        logger.Logger
	OrdersGroup   *echo.Group `name:"order-echo-group"`
	Validator     *validator.Validate
}
