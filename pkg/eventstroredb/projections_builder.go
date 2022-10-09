package eventstroredb

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/es/contracts/projection"
)

type ProjectionsBuilder interface {
	AddProjection(projection projection.IProjection) ProjectionsBuilder
	AddProjections(projections []projection.IProjection) ProjectionsBuilder
	Build() *ProjectionsConfigurations
}
