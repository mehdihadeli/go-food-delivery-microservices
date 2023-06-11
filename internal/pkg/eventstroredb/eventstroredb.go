package eventstroredb

import (
	"github.com/EventStore/EventStore-Client-Go/esdb"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/eventstroredb/config"
)

func NewEventStoreDB(cfg *config.EventStoreDbOptions) (*esdb.Client, error) {
	settings, err := esdb.ParseConnectionString(cfg.ConnectionString)
	if err != nil {
		return nil, err
	}

	return esdb.NewClient(settings)
}
