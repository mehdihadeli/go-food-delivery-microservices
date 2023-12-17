package types

type MessageEnvelope struct {
	Message IMessage
	Headers map[string]interface{}
}

func NewMessageEnvelope(
	message IMessage,
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

type MessageEnvelopeT[T IMessage] struct {
	*MessageEnvelope
	Message T
}

func NewMessageEnvelopeT[T IMessage](
	message T,
	headers map[string]interface{},
) *MessageEnvelopeT[T] {
	messageEnvelope := NewMessageEnvelope(message, headers)

	return &MessageEnvelopeT[T]{
		MessageEnvelope: messageEnvelope,
		Message:         message,
	}
}
