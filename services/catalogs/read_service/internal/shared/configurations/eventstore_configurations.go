package configurations

import (
	"fmt"
	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/es/store"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/eventstroredb"
)

func (ic *infrastructureConfigurator) configEventStore() (*esdb.Client, error, func()) {
	db, err := eventstroredb.NewEventStoreDB(ic.cfg.EventStoreConfig)
	if err != nil {
		return nil, err, nil
	}

	aggregateStore := store.NewAggregateStore(ic.log, db)
	fmt.Print(aggregateStore)

	return db, nil, func() {
		_ = db.Close() // nolint: errcheck
	}
}
