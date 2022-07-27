package store

import (
	"context"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/core/domain"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/es/contracts"
	uuid "github.com/satori/go.uuid"
)

// AggregateStore is responsible for loading and saving Aggregate.
type AggregateStore[T contracts.IEventSourcedAggregateRoot] interface {
	// Load loads the most recent version of an aggregate to provided  into params aggregate with a type and id.
	Load(ctx context.Context, aggregateId uuid.UUID) (T, error)

	// Store save the uncommitted events for an aggregate.
	Store(ctx context.Context, aggregate T, metadata *domain.Metadata) error

	// Exists check aggregate exists by AggregateId.
	Exists(ctx context.Context, aggregateId uuid.UUID) (bool, error)
}
