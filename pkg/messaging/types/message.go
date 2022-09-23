package types

import (
	"time"
)

type IMessage interface {
	GeMessageId() string
	GetCorrelationId() string
	GetCreated() time.Time
	GetEventTypeName() string
	SetCorrelationId(string)
	SetEventTypeName(string)
}

type Message struct {
	MessageId     string    `json:"messageId,omitempty"`
	CorrelationId string    `json:"correlationId"`
	Created       time.Time `json:"created"`
	EventType     string    `json:"eventType"`
}

func NewMessage(messageId string) *Message {
	return &Message{MessageId: messageId, Created: time.Now()}
}

func NewMessageWithTypeName(messageId string, eventTypeName string) *Message {
	return &Message{MessageId: messageId, Created: time.Now(), EventType: eventTypeName}
}

func (m *Message) GeMessageId() string {
	return m.MessageId
}

func (m *Message) GetCorrelationId() string {
	return m.CorrelationId
}

func (m *Message) GetCreated() time.Time {
	return m.Created
}

func (m *Message) SetCorrelationId(correlationId string) {
	m.CorrelationId = correlationId
}

func (m *Message) GetEventTypeName() string {
	return m.EventType
}

func (m *Message) SetEventTypeName(eventTypeName string) {
	m.EventType = eventTypeName
}
