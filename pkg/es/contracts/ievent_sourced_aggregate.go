package contracts

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/core/domain/contracts"
)

// IEventSourcedAggregateRoot contains all methods of AggregateBase
type IEventSourcedAggregateRoot interface {
	contracts.IAggregateRoot
	SetAppliedEvents(events []interface{})
	GetAppliedEvents() []interface{}
	RaiseEvent(event interface{}) error
	ToSnapshot()
	Load
	Apply
	When
}

// Apply process Aggregate Event
type Apply interface {
	Apply(event interface{}) error
}

// Load create Aggregate state from Event's.
type Load interface {
	Load(events []interface{}) error
}

type When interface {
	When(event interface{}) error
}
