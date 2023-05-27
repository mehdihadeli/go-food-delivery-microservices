package projections

import (
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/eventstroredb"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/messaging/bus"

	orderRepositories "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orders/internal/orders/data/repositories"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orders/internal/orders/projections"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orders/internal/shared/contracts"
)

func ConfigOrderProjections(builder eventstroredb.ProjectionsBuilder, infra *contracts.InfrastructureConfigurations, bus bus.Bus) {
	mongoOrderReadRepository := orderRepositories.NewMongoOrderReadRepository(infra.Log, infra.Cfg, infra.MongoClient)
	mongoOrderProjection := projections.NewMongoOrderProjection(mongoOrderReadRepository, bus, infra.Log)
	builder.AddProjection(mongoOrderProjection)

	elasticOrderReadRepository := orderRepositories.NewElasticOrderReadRepository(infra.Log, infra.Cfg, infra.ElasticClient)
	elasticOrderProjection := projections.NewElasticOrderProjection(elasticOrderReadRepository)
	builder.AddProjection(elasticOrderProjection)
}
