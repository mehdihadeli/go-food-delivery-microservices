package infrastructure

import (
	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/eventstroredb"
)

func (ic *infrastructureConfigurator) configEventStore() (*esdb.Client, error, func()) {
	db, err := eventstroredb.NewEventStoreDB(ic.cfg.EventStoreConfig)
	if err != nil {
		return nil, err, nil
	}

	return db, nil, func() {
		_ = db.Close() // nolint: errcheck
	}
}
