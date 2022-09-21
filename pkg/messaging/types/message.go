package types

import "time"

type IMessage interface {
	GeMessageId() string
	GetCorrelationId() string
	GetCreatedAt() time.Time
	SetCorrelationId(string)
}

type Message struct {
	MessageId     string    `json:"messageId,omitempty"`
	CorrelationId string    `json:"correlationId"`
	CreatedAt     time.Time `json:"createdAt,omitempty"`
}

func NewMessage(messageId string) *Message {
	return &Message{MessageId: messageId, CreatedAt: time.Now()}
}

func (m *Message) GeMessageId() string {
	return m.MessageId
}

func (m *Message) GetCorrelationId() string {
	return m.CorrelationId
}

func (m *Message) GetCreatedAt() time.Time {
	return m.CreatedAt
}

func (m *Message) SetCorrelationId(correlationId string) {
	m.CorrelationId = correlationId
}
