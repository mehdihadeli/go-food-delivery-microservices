package projections

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/eventstroredb"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/bus"
	orderRepositories "github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/data/repositories"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/projections"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/shared/contracts"
)

func ConfigOrderProjections(builder eventstroredb.ProjectionsBuilder, infra contracts.InfrastructureConfigurations, bus bus.Bus) {
	mongoOrderReadRepository := orderRepositories.NewMongoOrderReadRepository(infra.Log(), infra.Cfg(), infra.MongoClient())
	mongoOrderProjection := projections.NewMongoOrderProjection(mongoOrderReadRepository, bus, infra.Log())
	builder.AddProjection(mongoOrderProjection)

	elasticOrderReadRepository := orderRepositories.NewElasticOrderReadRepository(infra.Log(), infra.Cfg(), infra.ElasticClient())
	elasticOrderProjection := projections.NewElasticOrderProjection(elasticOrderReadRepository)
	builder.AddProjection(elasticOrderProjection)
}
