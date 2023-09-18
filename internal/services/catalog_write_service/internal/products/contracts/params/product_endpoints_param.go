package params

import (
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/web/route"

	"go.uber.org/fx"
)

type ProductsEndpointsParams struct {
	fx.In

	Endpoints []route.Endpoint `group:"product-routes"`
}
