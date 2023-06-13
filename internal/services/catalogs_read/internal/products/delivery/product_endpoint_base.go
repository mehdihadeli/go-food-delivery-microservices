package delivery

import (
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/read_service/internal/shared/contracts"
)

type ProductEndpointBase struct {
	ProductsGroup   *echo.Group
	CatalogsMetrics *contracts.CatalogsMetrics
	Logger          logger.Logger
	Validator       *validator.Validate
}

func NewProductEndpointBase(
	logger logger.Logger,
	validator *validator.Validate,
	productsGroup *echo.Group,
	catalogsMetrics *contracts.CatalogsMetrics,
) *ProductEndpointBase {
	return &ProductEndpointBase{
		ProductsGroup:   productsGroup,
		CatalogsMetrics: catalogsMetrics,
		Logger:          logger,
		Validator:       validator,
	}
}
