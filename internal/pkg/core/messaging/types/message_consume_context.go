package types

import (
	"time"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/core/metadata"
)

type MessageConsumeContext interface {
	MessageId() string
	CorrelationId() string
	MessageType() string
	Created() time.Time
	ContentType() string
	DeliveryTag() uint64
	Metadata() metadata.Metadata
	Message() IMessage
}

type messageConsumeContext struct {
	metadata      metadata.Metadata
	contentType   string
	messageType   string
	messageId     string
	created       time.Time
	tag           uint64
	correlationId string
	message       IMessage
}

func NewMessageConsumeContext(
	message IMessage,
	meta metadata.Metadata,
	contentType string,
	messageType string,
	created time.Time,
	deliveryTag uint64,
	messageId string,
	correlationId string,
) MessageConsumeContext {
	return &messageConsumeContext{
		message:       message,
		metadata:      meta,
		contentType:   contentType,
		messageId:     messageId,
		tag:           deliveryTag,
		created:       created,
		messageType:   messageType,
		correlationId: correlationId,
	}
}

func (m *messageConsumeContext) Message() IMessage {
	return m.message
}

func (m *messageConsumeContext) MessageId() string {
	return m.messageId
}

func (m *messageConsumeContext) CorrelationId() string {
	return m.correlationId
}

func (m *messageConsumeContext) MessageType() string {
	return m.messageType
}

func (m *messageConsumeContext) ContentType() string {
	return m.contentType
}

func (m *messageConsumeContext) Metadata() metadata.Metadata {
	return m.metadata
}

func (m *messageConsumeContext) Created() time.Time {
	return m.created
}

func (m *messageConsumeContext) DeliveryTag() uint64 {
	return m.tag
}
