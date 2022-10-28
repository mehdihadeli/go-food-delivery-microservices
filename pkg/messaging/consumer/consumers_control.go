package consumer

import (
	"context"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/types"
)

type BusControl interface {
	// Start starts all consumers
	Start(ctx context.Context) error
	// Stop stops all consumers
	Stop(ctx context.Context) error

	AddMessageConsumedHandler(func(message types.IMessage))
}
