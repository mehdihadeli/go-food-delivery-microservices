package configurations

import (
	"context"
	"fmt"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/mediatr"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/shared/configurations"
)

type ProductEndpointsConfigurator struct {
	*ProductModuleConfigurations
}

type ProductEndpointsConfigurations struct {
	*configurations.Infrastructure
	*mediatr.Mediator
}

func (pc *ProductEndpointsConfigurator) configEndpoints(ctx context.Context) {
	endpointsConfigurations := &ProductEndpointsConfigurations{Infrastructure: pc.Infrastructure, Mediator: pc.Mediator}
	fmt.Print(endpointsConfigurations)

	//createProductEndpoint := creating_product.NewCreteProductEndpoint(endpointsConfigurations)
	//createProductEndpoint.MapRoute()
}
