package configurations

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/mediatr"
)

func (c *productsModuleConfigurator) configEndpoints(ctx context.Context, mediator *mediatr.Mediator, group *echo.Group) {
	fmt.Print(c)

	//createProductEndpoint := creating_product.NewCreteProductEndpoint(productEndpointBase)
	//createProductEndpoint.MapRoute()
}
