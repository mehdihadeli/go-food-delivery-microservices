package projections

import (
	"fmt"
	orderRepositories "github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/data/repositories"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/projections"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/shared/configurations/infrastructure"
)

func ConfigOrderProjections(infra *infrastructure.InfrastructureConfiguration) {
	mongoOrderReadRepository := orderRepositories.NewMongoOrderReadRepository(infra.Log, infra.Cfg, infra.MongoClient)
	elasticOrderReadRepository := orderRepositories.NewElasticOrderReadRepository(infra.Log, infra.Cfg, infra.ElasticClient)

	mongoOrderProjection := projections.NewMongoOrderProjection(mongoOrderReadRepository, infra.Producer, infra.Log)
	infra.Projections = append(infra.Projections, mongoOrderProjection)

	elasticOrderProjection := projections.NewElasticOrderProjection(elasticOrderReadRepository)
	fmt.Println(elasticOrderProjection)
	//c.Projections = append(c.Projections, elasticOrderProjection)
}
