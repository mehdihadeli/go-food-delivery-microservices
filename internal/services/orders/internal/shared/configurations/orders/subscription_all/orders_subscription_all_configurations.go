package subscriptionAll

import (
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/eventstroredb"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/messaging/bus"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orders/internal/orders/configurations/projections"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orders/internal/shared/contracts"
)

func ConfigOrdersSubscriptionAllWorker(infra *contracts.InfrastructureConfigurations, bus bus.Bus) (eventstroredb.EsdbSubscriptionAllWorker, error) {
	esdbWorker := eventstroredb.NewEsdbSubscriptionAllWorker(
		infra.Log,
		infra.Esdb,
		infra.Cfg.EventStoreConfig,
		infra.EsdbSerializer,
		infra.CheckpointRepository,
		func(builder eventstroredb.ProjectionsBuilder) {
			projections.ConfigOrderProjections(builder, infra, bus)
		})

	return esdbWorker, nil
}
