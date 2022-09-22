package infrastructure

import (
	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/core/serializer/json"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/es/contracts"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/eventstroredb"
)

func (ic *infrastructureConfigurator) configEventStore() (*esdb.Client, contracts.SubscriptionCheckpointRepository, *eventstroredb.EsdbSerializer, error, func()) {
	db, err := eventstroredb.NewEventStoreDB(ic.cfg.EventStoreConfig)
	if err != nil {
		return nil, nil, nil, err, nil
	}

	esdbSerializer := eventstroredb.NewEsdbSerializer(json.NewJsonMetadataSerializer(), json.NewJsonEventSerializer())
	subscriptionRepository := eventstroredb.NewEsdbSubscriptionCheckpointRepository(db, ic.log, esdbSerializer)

	return db, subscriptionRepository, esdbSerializer, nil, func() {
		_ = db.Close() // nolint: errcheck
	}
}
