package eventstroredb

import (
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/eventstroredb/config"

	"github.com/EventStore/EventStore-Client-Go/esdb"
)

func NewEventStoreDB(cfg *config.EventStoreDbOptions) (*esdb.Client, error) {
	settings, err := esdb.ParseConnectionString(cfg.GrpcEndPoint())
	if err != nil {
		return nil, err
	}

	return esdb.NewClient(settings)
}
