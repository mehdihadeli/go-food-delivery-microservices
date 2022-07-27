package store

import (
	"context"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/es"
)

// EventStore is an interface for an Event sourcing event store.
type EventStore interface {
	// SaveEvents appends all events in the Event stream to the store.
	SaveEvents(ctx context.Context, streamID string, events []*es.ESEvent) error

	// LoadEvents loads all events for the Aggregate id from the store.
	LoadEvents(ctx context.Context, streamID string) ([]*es.ESEvent, error)
}
