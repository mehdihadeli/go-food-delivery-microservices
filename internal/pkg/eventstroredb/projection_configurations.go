package eventstroredb

import (
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/es/contracts/projection"
)

type ProjectionsConfigurations struct {
	Projections []projection.IProjection
}
