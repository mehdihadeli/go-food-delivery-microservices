package types

type IMessage interface {
	MessageId() string
	CorrelationId() string
	SetCorrelationId(string)
}

type Message struct {
	correlationId string
	messageId     string
}

func NewMessage(messageId string) *Message {
	return &Message{messageId: messageId}
}

func (m *Message) MessageId() string {
	return m.messageId
}

func (m *Message) CorrelationId() string {
	return m.correlationId
}

func (m *Message) SetCorrelationId(correlationId string) {
	m.correlationId = correlationId
}
