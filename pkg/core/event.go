package core

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

type IEvent interface {
	EventId() uuid.UUID
	EventType() string
	OccurredOn() time.Time
	EventVersion() int64
	SetEventVersion(version int64)
}

type Event struct {
	eventId      uuid.UUID
	eventType    string
	occurredOn   time.Time
	eventVersion int64
}

type EventDataModel struct {
	EventId      uuid.UUID `json:"eventId" bson:"eventId,omitempty"`
	EventType    string    `json:"eventType" bson:"eventType,omitempty"`
	OccurredOn   time.Time `json:"occurredOn" bson:"occurredOn,omitempty"`
	EventVersion int64     `json:"eventVersion" bson:"eventVersion,omitempty"`
}

func NewEvent(eventType string) *Event {
	return &Event{
		eventId:      uuid.NewV4(),
		occurredOn:   time.Now(),
		eventVersion: -1,
		eventType:    eventType,
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

func (e *Event) EventVersion() int64 {
	return e.eventVersion
}

func (e *Event) SetEventVersion(version int64) {
	e.eventVersion = version
}
