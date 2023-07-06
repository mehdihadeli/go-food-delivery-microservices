package eventstoredb

import (
	"context"
	"testing"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/eventstroredb/config"
)

var EventstoreDBContainerOptionsDecorator = func(t *testing.T, ctx context.Context) interface{} {
	return func(c *config.EventStoreDbOptions) (*config.EventStoreDbOptions, error) {
		newOption, err := NewEventstoreDBTestContainers().CreatingContainerOptions(ctx, t)
		if err != nil {
			return nil, err
		}
		newOption.Subscription = c.Subscription

		return newOption, nil
	}
}

var ReplaceEventStoreContainerOptions = func(t *testing.T, options *config.EventStoreDbOptions, ctx context.Context) error {
	newOption, err := NewEventstoreDBTestContainers().CreatingContainerOptions(ctx, t)
	if err != nil {
		return err
	}

	options.HttpPort = newOption.HttpPort
	options.TcpPort = newOption.TcpPort
	options.Host = newOption.Host

	return nil
}
