package core

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

type IEvent interface {
	EventId() uuid.UUID
	EventType() string
	OccurredOn() time.Time
}

type Event struct {
	eventId    uuid.UUID
	eventType  string
	occurredOn time.Time
}

func NewEvent(eventType string) *Event {
	return &Event{
		eventId:    uuid.NewV4(),
		occurredOn: time.Now(),
		eventType:  eventType,
	}
}

func (e *Event) EventId() uuid.UUID {
	return e.eventId
}

func (e *Event) EventType() string {
	return e.eventType
}

func (e *Event) OccurredOn() time.Time {
	return e.occurredOn
}
