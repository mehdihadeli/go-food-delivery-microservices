package contracts

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/core/domain/types"
	uuid "github.com/satori/go.uuid"
)

// IAggregateRoot contains all methods of AggregateBase
type IAggregateRoot interface {
	GetID() uuid.UUID
	SetID(id uuid.UUID)
	GetVersion() int64
	SetType(aggregateType types.AggregateType)
	GetType() types.AggregateType
	AddEvent(event interface{})
	MarkUncommittedEventAsCommitted()
	HasUncommittedEvents() bool
	GetUncommittedEvents() []interface{}
	String() string
}
