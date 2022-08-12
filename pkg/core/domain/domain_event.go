package domain

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/core"
	uuid "github.com/satori/go.uuid"
)

type IDomainEvent interface {
	core.IEvent
	GetAggregateId() uuid.UUID
	GetAggregateSequenceNumber() int64
	WithAggregate(aggregateId uuid.UUID, aggregateSequenceNumber int64) *DomainEvent
}

type DomainEvent struct {
	*core.Event
	aggregateId             uuid.UUID
	aggregateSequenceNumber int64
}

type DomainEventDataModel struct {
	*core.EventDataModel
	AggregateId             uuid.UUID `json:"aggregateId" bson:"aggregateId,omitempty"`
	AggregateSequenceNumber int64     `json:"aggregateSequenceNumber" bson:"aggregateSequenceNumber,omitempty"`
}

func NewDomainEvent(aggregateId uuid.UUID, aggregateSequenceNumber int64, eventType string) *DomainEvent {
	return &DomainEvent{
		Event:                   core.NewEvent(eventType),
		aggregateId:             aggregateId,
		aggregateSequenceNumber: aggregateSequenceNumber,
	}
}

func (d *DomainEvent) GetAggregateId() uuid.UUID {
	return d.aggregateId
}

func (d *DomainEvent) GetAggregateSequenceNumber() int64 {
	return d.aggregateSequenceNumber
}

func (d *DomainEvent) WithAggregate(aggregateId uuid.UUID, aggregateSequenceNumber int64) *DomainEvent {
	d.aggregateId = aggregateId
	d.aggregateSequenceNumber = aggregateSequenceNumber

	return d
}
