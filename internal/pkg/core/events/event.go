package events

import (
	"time"

	typeMapper "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/reflection/typemapper"

	uuid "github.com/satori/go.uuid"
)

type IEvent interface {
	GetEventId() uuid.UUID
	GetOccurredOn() time.Time
	// GetEventTypeName get short type name of the event - we use event short type name instead of full type name because this event in other receiver packages could have different package name
	GetEventTypeName() string
	GetEventFullTypeName() string
}

type Event struct {
	EventId    uuid.UUID `json:"event_id"`
	EventType  string    `json:"event_type"`
	OccurredOn time.Time `json:"occurred_on"`
}

func NewEvent(eventType string) *Event {
	return &Event{
		EventId:    uuid.NewV4(),
		OccurredOn: time.Now(),
		EventType:  eventType,
	}
}

func (e *Event) GetEventId() uuid.UUID {
	return e.EventId
}

func (e *Event) GetEventType() string {
	return e.EventType
}

func (e *Event) GetOccurredOn() time.Time {
	return e.OccurredOn
}

func (e *Event) GetEventTypeName() string {
	return typeMapper.GetTypeName(e)
}

func (e *Event) GetEventFullTypeName() string {
	return typeMapper.GetFullTypeName(e)
}

func IsEvent(obj interface{}) bool {
	if _, ok := obj.(IEvent); ok {
		return true
	}

	return false
}
