package params

import (
	"go.uber.org/fx"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/es/contracts/projection"
)

type OrderProjectionParams struct {
	fx.In

	Projections []projection.IProjection `group:"projections"`
}
