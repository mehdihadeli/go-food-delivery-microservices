package configurations

import (
	"context"
	"fmt"
)

func (pm *ProductModule) configEndpoints(ctx context.Context) {
	fmt.Print(pm)

	//createProductEndpoint := creating_product.NewCreteProductEndpoint(productEndpointBase)
	//createProductEndpoint.MapRoute()
}
