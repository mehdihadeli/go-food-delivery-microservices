package es

//https://www.eventstore.com/blog/what-is-event-sourcing
//https://www.eventstore.com/blog/event-sourcing-and-cqrs

import (
	"context"
	"fmt"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/core/domain/types"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/utils"
	uuid "github.com/satori/go.uuid"
)

const (
	newEventSourcesAggregateVersion      = -1 // used for init version in the EventStoreDB
	aggregateAppliedEventsInitialCap     = 10
	aggregateUncommittedEventsInitialCap = 10
)

// HandleCommand Aggregate commands' handler method
// Example
//
// func (a *OrderAggregate) HandleCommand(command interface{}) error {
//	switch c := command.(type) {
//	case *CreateOrderCommand:
//		return a.handleCreateOrderCommand(c)
//	case *OrderPaidCommand:
//		return a.handleOrderPaidCommand(c)
//	case *SubmitOrderCommand:
//		return a.handleSubmitOrderCommand(c)
//	default:
//		return errors.New("invalid command type")
//	}
//}
type HandleCommand interface {
	HandleCommand(ctx context.Context, command Command) error
}

type when func(event interface{}) error

// EventSourcedAggregateRoot base aggregate contains all main necessary fields
type EventSourcedAggregateRoot struct {
	ID                uuid.UUID           `json:"id" bson:"_id,omitempty"`
	Version           int64               `json:"version" bson:"version"`
	Type              types.AggregateType `json:"type" bson:"type"`
	UncommittedEvents []interface{}
	AppliedEvents     []interface{}
	withAppliedEvents bool
	when              when
}

func NewEventSourcedAggregateRoot(when when) *EventSourcedAggregateRoot {
	if when == nil {
		return nil
	}

	return &EventSourcedAggregateRoot{

		Version:           newEventSourcesAggregateVersion,
		UncommittedEvents: make([]interface{}, 0, aggregateUncommittedEventsInitialCap),
		AppliedEvents:     make([]interface{}, 0, aggregateAppliedEventsInitialCap),
		when:              when,
		withAppliedEvents: false,
	}
}

// SetID set EventSourcedAggregateRoot ID
func (a *EventSourcedAggregateRoot) SetID(id uuid.UUID) {
	a.ID = id
}

// GetID get EventSourcedAggregateRoot ID
func (a *EventSourcedAggregateRoot) GetID() uuid.UUID {
	return a.ID
}

// SetType set EventSourcedAggregateRoot AggregateType
func (a *EventSourcedAggregateRoot) SetType(aggregateType types.AggregateType) {
	a.Type = aggregateType
}

// GetType get EventSourcedAggregateRoot AggregateType
func (a *EventSourcedAggregateRoot) GetType() types.AggregateType {
	return a.Type
}

// GetAppliedEvents get EventSourcedAggregateRoot applied Event's
func (a *EventSourcedAggregateRoot) GetAppliedEvents() []interface{} {
	return a.AppliedEvents
}

// SetAppliedEvents set EventSourcedAggregateRoot applied Event's
func (a *EventSourcedAggregateRoot) SetAppliedEvents(events []interface{}) {
	a.AppliedEvents = events
}

// GetVersion get EventSourcedAggregateRoot version
func (a *EventSourcedAggregateRoot) GetVersion() int64 {
	return a.Version
}

// AddEvent add a new event to the EventSourcedAggregateRoot uncommitted Event's
func (a *EventSourcedAggregateRoot) AddEvent(event interface{}) {
	if utils.ContainsFunc(a.UncommittedEvents, func(e interface{}) bool {
		return e == event
	}) {
		return
	}
	a.UncommittedEvents = append(a.UncommittedEvents, event)
	a.Version++
}

// MarkUncommittedEventAsCommitted clear EventSourcedAggregateRoot uncommitted Event's
func (a *EventSourcedAggregateRoot) MarkUncommittedEventAsCommitted() {
	a.UncommittedEvents = make([]interface{}, 0, aggregateUncommittedEventsInitialCap)
}

// HasUncommittedEvents returns true if EventSourcedAggregateRoot has uncommitted Event's
func (a *EventSourcedAggregateRoot) HasUncommittedEvents() bool {
	return len(a.UncommittedEvents) > 0
}

// GetUncommittedEvents get EventSourcedAggregateRoot uncommitted Event's
func (a *EventSourcedAggregateRoot) GetUncommittedEvents() []interface{} {
	return a.UncommittedEvents
}

// Load add existing events from event store to aggregate using When interface method
func (a *EventSourcedAggregateRoot) Load(events []interface{}) error {

	for _, evt := range events {
		if err := a.when(evt); err != nil {
			return err
		}

		if a.withAppliedEvents {
			a.AppliedEvents = append(a.AppliedEvents, evt)
		}
		a.Version++
	}

	return nil
}

// Ref: https://www.eventstore.com/blog/what-is-event-sourcing

// Apply push event to aggregate uncommitted events using When method
func (a *EventSourcedAggregateRoot) Apply(event interface{}) error {

	if err := a.when(event); err != nil {
		return err
	}

	a.AddEvent(event)

	return nil
}

// RaiseEvent push event to aggregate applied events using When method, used for load directly from eventstore
func (a *EventSourcedAggregateRoot) RaiseEvent(event interface{}) error {

	if err := a.when(event); err != nil {
		return err
	}

	if a.withAppliedEvents {
		a.AppliedEvents = append(a.AppliedEvents, event)
	}

	a.Version++

	return nil
}

// ToSnapshot prepare EventSourcedAggregateRoot for saving Snapshot.
func (a *EventSourcedAggregateRoot) ToSnapshot() {
	if a.withAppliedEvents {
		a.AppliedEvents = append(a.AppliedEvents, a.UncommittedEvents...)
	}
	a.MarkUncommittedEventAsCommitted()
}

func (a *EventSourcedAggregateRoot) String() string {
	return fmt.Sprintf("ID: {%s}, Version: {%v}, Type: {%v}, AppliedEvents: {%v}, UncommittedEvents: {%v}",
		a.GetID(),
		a.GetVersion(),
		a.GetType(),
		len(a.GetAppliedEvents()),
		len(a.GetUncommittedEvents()),
	)
}
