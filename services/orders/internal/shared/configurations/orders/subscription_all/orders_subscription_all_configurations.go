package subscriptionAll

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/eventstroredb"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/bus"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/configurations/projections"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/shared/contracts"
)

func ConfigOrdersSubscriptionAllWorker(infra contracts.InfrastructureConfigurations, bus bus.Bus) (eventstroredb.EsdbSubscriptionAllWorker, error) {
	esdbWorker := eventstroredb.NewEsdbSubscriptionAllWorker(
		infra.Log(),
		infra.Esdb(),
		infra.Cfg().EventStoreConfig,
		infra.EsdbSerializer(),
		infra.CheckpointRepository(),
		func(builder eventstroredb.ProjectionsBuilder) {
			projections.ConfigOrderProjections(builder, infra, bus)
		})

	return esdbWorker, nil
}
