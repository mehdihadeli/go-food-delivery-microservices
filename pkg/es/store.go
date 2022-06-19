package es

import (
	"context"
)

// AggregateStore is responsible for loading and saving aggregates.
type AggregateStore interface {
	// Load loads the most recent version of an aggregate to provided  into params aggregate with a type and id.
	Load(ctx context.Context, aggregate Aggregate) error

	// Save saves the uncommitted events for an aggregate.
	Save(ctx context.Context, aggregate Aggregate) error

	// Exists check aggregate exists by id.
	Exists(ctx context.Context, streamID string) error

	//EventStore
	//SnapshotStore
}

// EventStore is an interface for an event sourcing event store.
type EventStore interface {
	// SaveEvents appends all events in the event stream to the store.
	SaveEvents(ctx context.Context, streamID string, events []Event) error

	// LoadEvents loads all events for the aggregate id from the store.
	LoadEvents(ctx context.Context, streamID string) ([]Event, error)
}

// SnapshotStore is an interface for an event sourcing snapshot store.
type SnapshotStore interface {
	// SaveSnapshot save aggregate snapshot.
	SaveSnapshot(ctx context.Context, aggregate Aggregate) error

	// GetSnapshot load aggregate snapshot.
	GetSnapshot(ctx context.Context, id string) (*Snapshot, error)
}
