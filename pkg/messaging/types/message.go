package types

type IMessage interface {
	GeMessageId() string
	GetCorrelationId() string
	SetCorrelationId(string)
}

type Message struct {
	CorrelationId string
	MessageId     string
}

func NewMessage(messageId string) *Message {
	return &Message{MessageId: messageId}
}

func (m *Message) GeMessageId() string {
	return m.MessageId
}

func (m *Message) GetCorrelationId() string {
	return m.CorrelationId
}

func (m *Message) SetCorrelationId(correlationId string) {
	m.CorrelationId = correlationId
}
