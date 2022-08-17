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
	AggregateId             uuid.UUID
	AggregateSequenceNumber int64
}

func NewDomainEvent(aggregateId uuid.UUID, aggregateSequenceNumber int64, eventType string) *DomainEvent {
	domainEvent := &DomainEvent{
		Event:                   core.NewEvent(eventType),
		AggregateId:             aggregateId,
		AggregateSequenceNumber: aggregateSequenceNumber,
	}
	domainEvent.Event = core.NewEvent(eventType)

	return domainEvent
}

func (d *DomainEvent) GetAggregateId() uuid.UUID {
	return d.AggregateId
}

func (d *DomainEvent) GetAggregateSequenceNumber() int64 {
	return d.AggregateSequenceNumber
}

func (d *DomainEvent) WithAggregate(aggregateId uuid.UUID, aggregateSequenceNumber int64) *DomainEvent {
	d.AggregateId = aggregateId
	d.AggregateSequenceNumber = aggregateSequenceNumber

	return d
}
