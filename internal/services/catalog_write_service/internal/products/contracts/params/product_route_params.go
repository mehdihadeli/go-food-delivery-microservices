package params

import (
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
	"go.uber.org/fx"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/shared/contracts"
)

type ProductRouteParams struct {
	fx.In

	CatalogsMetrics *contracts.CatalogsMetrics
	Logger          logger.Logger
	ProductsGroup   *echo.Group `name:"product-echo-group"`
	Validator       *validator.Validate
}
