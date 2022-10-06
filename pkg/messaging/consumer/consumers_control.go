package consumer

import "context"

type ConsumersControl interface {
	// Start starts all consumers
	Start(ctx context.Context) error
	// Stop stops all consumers
	Stop(ctx context.Context) error
}
