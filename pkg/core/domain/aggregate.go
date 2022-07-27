package domain

import (
	"fmt"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/core/domain/types"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/utils"
	uuid "github.com/satori/go.uuid"
)

const (
	newAggregateVersion                  = 0
	aggregateUncommittedEventsInitialCap = 10
)

// AggregateRoot base aggregate contains all main necessary fields
type AggregateRoot struct {
	ID                uuid.UUID
	Version           int64
	UncommittedEvents []interface{}
	Type              types.AggregateType
}

func NewAggregateRoot() *AggregateRoot {

	return &AggregateRoot{
		Version:           newAggregateVersion,
		UncommittedEvents: make([]interface{}, 0, aggregateUncommittedEventsInitialCap),
	}
}

// SetID set AggregateRoot ID
func (a *AggregateRoot) SetID(id uuid.UUID) {
	a.ID = id
}

// GetID get AggregateRoot ID
func (a *AggregateRoot) GetID() uuid.UUID {
	return a.ID
}

// SetType set AggregateRoot AggregateType
func (a *AggregateRoot) SetType(aggregateType types.AggregateType) {
	a.Type = aggregateType
}

// GetType get AggregateRoot AggregateType
func (a *AggregateRoot) GetType() types.AggregateType {
	return a.Type
}

// GetVersion get AggregateRoot version
func (a *AggregateRoot) GetVersion() int64 {
	return a.Version
}

// AddEvent add a new event to the AggregateRoot uncommitted Event's
func (a *AggregateRoot) AddEvent(event interface{}) {
	if utils.ContainsFunc(a.UncommittedEvents, func(e interface{}) bool {
		return e == event
	}) {
		return
	}
	a.UncommittedEvents = append(a.UncommittedEvents, event)
}

// MarkUncommittedEventAsCommitted clear AggregateRoot uncommitted Event's
func (a *AggregateRoot) MarkUncommittedEventAsCommitted() {
	a.UncommittedEvents = make([]interface{}, 0, aggregateUncommittedEventsInitialCap)
}

// HasUncommittedEvents returns true if AggregateRoot has uncommitted Event's
func (a *AggregateRoot) HasUncommittedEvents() bool {
	return len(a.UncommittedEvents) > 0
}

// GetUncommittedEvents get AggregateRoot uncommitted Event's
func (a *AggregateRoot) GetUncommittedEvents() []interface{} {
	return a.UncommittedEvents
}

func (a *AggregateRoot) String() string {
	return fmt.Sprintf("ID: {%s}, Version: {%v}, Type: {%v} , UncommittedEvents: {%v}",
		a.GetID(),
		a.GetVersion(),
		a.GetType(),
		len(a.GetUncommittedEvents()),
	)
}
