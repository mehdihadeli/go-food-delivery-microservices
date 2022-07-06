package configurations

import (
	"context"
	"fmt"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/mediatr"
)

func (c *productsModuleConfigurator) configGrpc(ctx context.Context, mediator *mediatr.Mediator) {
	fmt.Print(c)
}
