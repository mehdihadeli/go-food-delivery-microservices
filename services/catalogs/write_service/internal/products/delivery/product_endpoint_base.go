package delivery

import (
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/mediatr"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/contracts"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/shared/configurations"
)

type ProductEndpointBase struct {
	*configurations.Infrastructure
	Mediator          *mediatr.Mediator
	ProductRepository contracts.ProductRepository
	ProductsGroup     *echo.Group
}
