package eventstroredb

import (
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/es/contracts/projection"
)

type ProjectionsConfigurations struct {
	Projections []projection.IProjection
}
