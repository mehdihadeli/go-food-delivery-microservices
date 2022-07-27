package es

import (
	"encoding/json"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/core/domain/types"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/es/contracts"
	uuid "github.com/satori/go.uuid"
)

// Snapshot Event Sourcing Snapshotting is an optimisation that reduces time spent on reading event from an event store.
type Snapshot struct {
	ID      uuid.UUID           `json:"id"`
	Type    types.AggregateType `json:"type"`
	State   []byte              `json:"state"`
	Version uint64              `json:"version"`
}

// NewSnapshotFromAggregate create new snapshot from the Aggregate state.
func NewSnapshotFromAggregate(aggregate contracts.IEventSourcedAggregateRoot) (*Snapshot, error) {

	aggregateBytes, err := json.Marshal(aggregate)
	if err != nil {
		return nil, err
	}

	return &Snapshot{
		ID:      aggregate.GetID(),
		Type:    aggregate.GetType(),
		State:   aggregateBytes,
		Version: uint64(aggregate.GetVersion()),
	}, nil
}
