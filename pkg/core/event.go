package core

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

type IEvent interface {
	GetEventId() uuid.UUID
	GetEventType() string
	GetOccurredOn() time.Time
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
