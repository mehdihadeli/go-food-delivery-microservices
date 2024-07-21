package eventstoredb

import (
	"context"
	"testing"

	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/eventstroredb/config"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/logger"
)

var EventstoreDBContainerOptionsDecorator = func(t *testing.T, ctx context.Context) interface{} {
	return func(c *config.EventStoreDbOptions, logger logger.Logger) (*config.EventStoreDbOptions, error) {
		newOption, err := NewEventstoreDBTestContainers(logger).PopulateContainerOptions(ctx, t)
		if err != nil {
			return nil, err
		}
		newOption.Subscription = c.Subscription

		return newOption, nil
	}
}

var ReplaceEventStoreContainerOptions = func(t *testing.T, options *config.EventStoreDbOptions, ctx context.Context, logger logger.Logger) error {
	newOption, err := NewEventstoreDBTestContainers(logger).PopulateContainerOptions(ctx, t)
	if err != nil {
		return err
	}

	options.HttpPort = newOption.HttpPort
	options.TcpPort = newOption.TcpPort
	options.Host = newOption.Host

	return nil
}
