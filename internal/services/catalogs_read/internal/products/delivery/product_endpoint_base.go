package delivery

import (
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/messaging/bus"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/read_service/internal/shared/contracts"
)

type ProductEndpointBase struct {
	*contracts.InfrastructureConfigurations
	ProductsGroup   *echo.Group
	CatalogsMetrics *contracts.CatalogsMetrics
	Bus             bus.Bus
}

func NewProductEndpointBase(infra *contracts.InfrastructureConfigurations, productsGroup *echo.Group, bus bus.Bus, catalogsMetrics *contracts.CatalogsMetrics) *ProductEndpointBase {
	return &ProductEndpointBase{InfrastructureConfigurations: infra, ProductsGroup: productsGroup, Bus: bus, CatalogsMetrics: catalogsMetrics}
}
