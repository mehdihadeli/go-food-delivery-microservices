package eventstroredb

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/es/contracts/projection"
)

type ProjectionsConfigurations struct {
	Projections []projection.IProjection
}
