package params

import (
	"go.uber.org/fx"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/web/route"
)

type ProductsEndpointsParams struct {
	fx.In

	Endpoints []route.Endpoint `group:"product-routes"`
}
