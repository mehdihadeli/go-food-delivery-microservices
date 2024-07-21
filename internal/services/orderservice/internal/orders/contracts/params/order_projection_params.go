package params

import (
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/es/contracts/projection"

	"go.uber.org/fx"
)

type OrderProjectionParams struct {
	fx.In

	Projections []projection.IProjection `group:"projections"`
}
