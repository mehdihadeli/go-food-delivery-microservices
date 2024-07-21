package domain

import (
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/core/events"
	expectedStreamVersion "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/es/models/stream_version"

	uuid "github.com/satori/go.uuid"
)

type IDomainEvent interface {
	events.IEvent
	GetAggregateId() uuid.UUID
	GetAggregateSequenceNumber() int64
	WithAggregate(aggregateId uuid.UUID, aggregateSequenceNumber int64) *DomainEvent
}

type DomainEvent struct {
	*events.Event
	AggregateId             uuid.UUID `json:"aggregate_id"`
	AggregateSequenceNumber int64     `json:"aggregate_sequence_number"`
}

func NewDomainEvent(eventType string) *DomainEvent {
	domainEvent := &DomainEvent{
		Event:                   events.NewEvent(eventType),
		AggregateSequenceNumber: expectedStreamVersion.NoStream.Value(),
	}
	domainEvent.Event = events.NewEvent(eventType)

	return domainEvent
}

func (d *DomainEvent) GetAggregateId() uuid.UUID {
	return d.AggregateId
}

func (d *DomainEvent) GetAggregateSequenceNumber() int64 {
	return d.AggregateSequenceNumber
}

func (d *DomainEvent) WithAggregate(
	aggregateId uuid.UUID,
	aggregateSequenceNumber int64,
) *DomainEvent {
	d.AggregateId = aggregateId
	d.AggregateSequenceNumber = aggregateSequenceNumber

	return d
}
