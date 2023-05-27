package delivery

import (
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/messaging/bus"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/write_service/internal/shared/contracts"
)

type ProductEndpointBase struct {
	*contracts.InfrastructureConfigurations
	ProductsGroup   *echo.Group
	CatalogsMetrics *contracts.CatalogsMetrics
	Bus             bus.Bus
}

func NewProductEndpointBase(infra *contracts.InfrastructureConfigurations, productsGroup *echo.Group, bus bus.Bus, catalogsMetrics *contracts.CatalogsMetrics) *ProductEndpointBase {
	return &ProductEndpointBase{ProductsGroup: productsGroup, InfrastructureConfigurations: infra, Bus: bus, CatalogsMetrics: catalogsMetrics}
}
