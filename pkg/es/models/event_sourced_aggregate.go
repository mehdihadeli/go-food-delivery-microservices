package models

//https://www.eventstore.com/blog/what-is-event-sourcing
//https://www.eventstore.com/blog/event-sourcing-and-cqrs

import (
	"emperror.dev/errors"
	"fmt"
	"github.com/ahmetb/go-linq/v3"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/core"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/core/domain"
	errors2 "github.com/mehdihadeli/store-golang-microservice-sample/pkg/es/errors"
	expectedStreamVersion "github.com/mehdihadeli/store-golang-microservice-sample/pkg/es/models/stream_version"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/serializer/jsonSerializer"
	uuid "github.com/satori/go.uuid"
)

type WhenFunc func(event domain.IDomainEvent) error

type When interface {
	// When Update the aggregate state with new events that are added to the event store and also for events that are already in the event store without increasing the version.
	When(event domain.IDomainEvent) error
}

type fold interface {
	// Restore the aggregate state with events that are loaded form the event store and increase the current version and last commit version.
	fold(event domain.IDomainEvent, metadata *core.Metadata) error
}

type Apply interface {
	// Apply a new event to the aggregate state, adds the event to the list of pending changes,
	// and increases the `CurrentVersion` property and `LastCommittedVersion` will be unchanged.
	Apply(event domain.IDomainEvent, isNew bool) error
}

type AggregateStateProjection interface {
	Apply
	fold
}

// IHaveEventSourcedAggregate this interface should implement by actual aggregate root class in our domain
type IHaveEventSourcedAggregate interface {
	When
	NewEmptyAggregate()
	IEventSourcedAggregateRoot
}

// IEventSourcedAggregateRoot contains all methods of AggregateBase
type IEventSourcedAggregateRoot interface {
	domain.IEntity

	// OriginalVersion Gets the original version is the aggregate version we got from the store. This is used to ensure optimistic concurrency,
	// to check if there were no changes made to the aggregate state between load and save for the current operation.
	OriginalVersion() int64

	SetOriginalVersion(version int64)

	// CurrentVersion Gets the current version is set to original version when the aggregate is loaded from the store.
	// It should increase for each state transition performed within the scope of the current operation.
	CurrentVersion() int64

	// AddDomainEvents adds a new domain event to the aggregate's uncommitted events.
	AddDomainEvents(event domain.IDomainEvent) error

	// MarkUncommittedEventAsCommitted Mark all changes (events) as committed, clears uncommitted changes and updates the current version of the aggregate.
	MarkUncommittedEventAsCommitted()

	// HasUncommittedEvents Does the aggregate have change that have not been committed to storage
	HasUncommittedEvents() bool

	// UncommittedEvents Gets a list of uncommitted events for this aggregate.
	UncommittedEvents() []domain.IDomainEvent

	// LoadFromHistory Loads the current state of the aggregate from a list of events.
	LoadFromHistory(events []domain.IDomainEvent, metadata *core.Metadata) error

	AggregateStateProjection
}

// EventSourcedAggregateRoot base aggregate contains all main necessary fields
type EventSourcedAggregateRoot struct {
	*domain.Entity
	originalVersion   int64
	currentVersion    int64
	uncommittedEvents []domain.IDomainEvent
	when              WhenFunc
}

type EventSourcedAggregateRootDataModel struct {
	*domain.EntityDataModel
	OriginalVersion int64 `json:"originalVersion" bson:"originalVersion"`
}

func NewEventSourcedAggregateRootWithId(id uuid.UUID, aggregateType string, when WhenFunc) *EventSourcedAggregateRoot {
	if when == nil {
		return nil
	}

	aggregate := &EventSourcedAggregateRoot{
		originalVersion: expectedStreamVersion.NoStream.Value(),
		currentVersion:  expectedStreamVersion.NoStream.Value(),
		when:            when,
	}

	aggregate.Entity = domain.NewEntityWithId(id, aggregateType)

	return aggregate
}

func NewEventSourcedAggregateRoot(aggregateType string, when WhenFunc) *EventSourcedAggregateRoot {
	if when == nil {
		return nil
	}

	aggregate := &EventSourcedAggregateRoot{
		originalVersion: expectedStreamVersion.NoStream.Value(),
		currentVersion:  expectedStreamVersion.NoStream.Value(),
		when:            when,
	}

	aggregate.Entity = domain.NewEntity(aggregateType)

	return aggregate
}

func (a *EventSourcedAggregateRoot) OriginalVersion() int64 {
	return a.originalVersion
}

func (a *EventSourcedAggregateRoot) SetOriginalVersion(version int64) {
	a.originalVersion = version
}

func (a *EventSourcedAggregateRoot) CurrentVersion() int64 {
	return a.currentVersion
}

func (a *EventSourcedAggregateRoot) AddDomainEvents(event domain.IDomainEvent) error {
	exists := linq.From(a.uncommittedEvents).AnyWithT(func(e domain.IDomainEvent) bool {
		return e.GetEventId() == event.GetEventId()
	})

	if exists {
		return errors2.EventAlreadyExistsError
	}
	event.WithAggregate(a.Id(), a.CurrentVersion()+1)
	a.uncommittedEvents = append(a.uncommittedEvents, event)

	return nil
}

func (a *EventSourcedAggregateRoot) MarkUncommittedEventAsCommitted() {
	a.uncommittedEvents = nil
}

func (a *EventSourcedAggregateRoot) HasUncommittedEvents() bool {
	return len(a.uncommittedEvents) > 0
}

func (a *EventSourcedAggregateRoot) UncommittedEvents() []domain.IDomainEvent {
	return a.uncommittedEvents
}

func (a *EventSourcedAggregateRoot) LoadFromHistory(events []domain.IDomainEvent, metadata *core.Metadata) error {
	for _, event := range events {
		err := a.fold(event, metadata)
		if err != nil {
			return errors.WrapIf(err, "[EventSourcedAggregateRoot_LoadFromHistory:fold] error in loading event from history")
		}
	}

	return nil
}

func (a *EventSourcedAggregateRoot) Apply(event domain.IDomainEvent, isNew bool) error {
	if isNew {
		err := a.AddDomainEvents(event)
		if err != nil {
			return errors.WrapIf(err, "[EventSourcedAggregateRoot_Apply:AddDomainEvents] error in adding domain event to the domain events list")
		}
	}
	err := a.when(event)
	if err != nil {
		return errors.WrapIf(err, "[EventSourcedAggregateRoot_Apply:when] error in the whenFunc")
	}
	a.currentVersion++

	return nil
}

func (a *EventSourcedAggregateRoot) fold(event domain.IDomainEvent, metadata *core.Metadata) error {
	err := a.when(event)
	if err != nil {
		return errors.WrapIf(err, "[EventSourcedAggregateRoot_fold:when] error in the applying whenFunc")
	}
	a.originalVersion++
	a.currentVersion++

	return nil
}

func (a *EventSourcedAggregateRoot) String() string {
	return fmt.Sprintf("Aggregate json is: %s", jsonSerializer.ColoredPrettyPrint(a))
}
