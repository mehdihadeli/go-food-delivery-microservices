package types

type MessageEnvelope struct {
	Message interface{}
	Headers map[string]interface{}
}

func NewMessageEnvelope(
	message interface{},
	headers map[string]interface{},
) *MessageEnvelope {
	if headers == nil {
		headers = make(map[string]interface{})
	}

	return &MessageEnvelope{
		Message: message,
		Headers: headers,
	}
}

type MessageEnvelopeTMessage struct {
	*MessageEnvelope
	MessageTMessage interface{}
}

func NewMessageEnvelopeTMessage(
	messageTMessage interface{},
	headers map[string]interface{},
) *MessageEnvelopeTMessage {
	messageEnvelope := NewMessageEnvelope(messageTMessage, headers)

	return &MessageEnvelopeTMessage{
		MessageEnvelope: messageEnvelope,
		MessageTMessage: messageTMessage,
	}
}
