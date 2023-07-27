package domain

import (
	"fmt"

	"emperror.dev/errors"
	"github.com/ahmetb/go-linq/v3"
	"github.com/goccy/go-json"
	uuid "github.com/satori/go.uuid"
)

const (
	newAggregateVersion = 0
)

// AggregateRoot base aggregate contains all main necessary fields
type AggregateRoot struct {
	*Entity
	originalVersion   int64
	uncommittedEvents []IDomainEvent
}

type AggregateDataModel struct {
	*EntityDataModel
	OriginalVersion int64 `json:"originalVersion" bson:"originalVersion,omitempty"`
}

type IAggregateRoot interface {
	IEntity

	// OriginalVersion Gets the original version is the aggregate version we got from the store. This is used to ensure optimistic concurrency,
	// to check if there were no changes made to the aggregate state between load and save for the current operation.
	OriginalVersion() int64

	// AddDomainEvents adds a new domain_events event to the aggregate's uncommitted events.
	AddDomainEvents(event IDomainEvent) error

	// MarkUncommittedEventAsCommitted Mark all changes (events) as committed, clears uncommitted changes and updates the current version of the aggregate.
	MarkUncommittedEventAsCommitted()

	// HasUncommittedEvents Does the aggregate have change that have not been committed to storage
	HasUncommittedEvents() bool

	// GetUncommittedEvents Gets a list of uncommitted events for this aggregate.
	GetUncommittedEvents() []IDomainEvent

	// String Returns a string representation of the aggregate.
	String() string
}

func NewAggregateRootWithId(id uuid.UUID, aggregateType string) *AggregateRoot {
	aggregate := &AggregateRoot{
		originalVersion: newAggregateVersion,
	}
	aggregate.Entity = NewEntityWithId(id, aggregateType)

	return aggregate
}

func NewAggregateRoot(aggregateType string) *AggregateRoot {
	aggregate := &AggregateRoot{
		originalVersion: newAggregateVersion,
	}
	aggregate.Entity = NewEntity(aggregateType)

	return aggregate
}

func (a *AggregateRoot) AddDomainEvent(event IDomainEvent) error {
	exists := linq.From(a.uncommittedEvents).Contains(event)
	if exists {
		return errors.New("event already exists")
	}
	a.uncommittedEvents = append(a.uncommittedEvents, event)

	return nil
}

func (a *AggregateRoot) OriginalVersion() int64 {
	return a.originalVersion
}

func (a *AggregateRoot) AddDomainEvents(event IDomainEvent) {
	if linq.From(a.uncommittedEvents).Contains(event) {
		return
	}

	a.uncommittedEvents = append(a.uncommittedEvents, event)
}

// MarkUncommittedEventAsCommitted clear AggregateRoot uncommitted domain_events events
func (a *AggregateRoot) MarkUncommittedEventAsCommitted() {
	a.uncommittedEvents = nil
}

// HasUncommittedEvents returns true if AggregateRoot has uncommitted domain_events events
func (a *AggregateRoot) HasUncommittedEvents() bool {
	return len(a.uncommittedEvents) > 0
}

// GetUncommittedEvents get AggregateRoot uncommitted domain_events events
func (a *AggregateRoot) GetUncommittedEvents() []IDomainEvent {
	return a.uncommittedEvents
}

func (a *AggregateRoot) String() string {
	j, _ := json.Marshal(a)
	return fmt.Sprintf("Aggregate json: %s", string(j))
}
