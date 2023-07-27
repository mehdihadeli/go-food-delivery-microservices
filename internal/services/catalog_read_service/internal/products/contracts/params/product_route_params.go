package params

import (
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogreadservice/internal/shared/contracts"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
)

type ProductRouteParams struct {
	fx.In

	CatalogsMetrics *contracts.CatalogsMetrics
	Logger          logger.Logger
	ProductsGroup   *echo.Group `name:"product-echo-group"`
	Validator       *validator.Validate
}
