package types

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/core/metadata"
	"time"
)

type MessageConsumeContextBase interface {
	MessageId() string
	CorrelationId() string
	MessageType() string
	Created() time.Time
	ContentType() string
	DeliveryTag() uint64
	Body() interface{}
	Metadata() metadata.Metadata
}

type MessageConsumeContextT[T IMessage] interface {
	MessageConsumeContextBase
	Message() T
	ToMessageConsumeContext() MessageConsumeContext
}

type MessageConsumeContext interface {
	MessageConsumeContextBase
	Message() IMessage
}

type messageConsumeContextBase struct {
	metadata      metadata.Metadata
	body          interface{}
	contentType   string
	messageType   string
	messageId     string
	created       time.Time
	tag           uint64
	correlationId string
}

type messageConsumeContextT[T IMessage] struct {
	MessageConsumeContextBase
	message T
}

type messageConsumeContext struct {
	MessageConsumeContextBase
	message IMessage
}

func NewMessageContextBase(body interface{}, meta metadata.Metadata, contentType string, messageType string, created time.Time, deliveryTag uint64, messageId string, correlationId string) MessageConsumeContextBase {
	return &messageConsumeContextBase{
		metadata:      meta,
		body:          body,
		contentType:   contentType,
		messageId:     messageId,
		tag:           deliveryTag,
		created:       created,
		messageType:   messageType,
		correlationId: correlationId,
	}
}
func NewMessageConsumeContextT[T IMessage](message T, meta metadata.Metadata, contentType string, messageType string, created time.Time, deliveryTag uint64, messageId string, correlationId string) MessageConsumeContextT[T] {
	return &messageConsumeContextT[T]{
		message:                   message,
		MessageConsumeContextBase: NewMessageContextBase(message, meta, contentType, messageType, created, deliveryTag, messageId, correlationId),
	}
}

func NewMessageConsumeContext(message IMessage, meta metadata.Metadata, contentType string, messageType string, created time.Time, deliveryTag uint64, messageId string, correlationId string) MessageConsumeContext {
	return &messageConsumeContext{
		message:                   message,
		MessageConsumeContextBase: NewMessageContextBase(message, meta, contentType, messageType, created, deliveryTag, messageId, correlationId),
	}
}

func (m *messageConsumeContext) Message() IMessage {
	return m.message
}

func (m *messageConsumeContextT[T]) Message() T {
	return m.message
}

func (m *messageConsumeContextT[T]) ToMessageConsumeContext() MessageConsumeContext {
	return NewMessageConsumeContext(m.Message(), m.Metadata(), m.ContentType(), m.MessageType(), m.Created(), m.DeliveryTag(), m.MessageId(), m.CorrelationId())
}

func (m *messageConsumeContextBase) MessageId() string {
	return m.messageId
}

func (m *messageConsumeContextBase) Body() interface{} {
	return m.body
}

func (m *messageConsumeContextBase) CorrelationId() string {
	return m.correlationId
}

func (m *messageConsumeContextBase) MessageType() string {
	return m.messageType
}

func (m *messageConsumeContextBase) ContentType() string {
	return m.contentType
}

func (m *messageConsumeContextBase) Metadata() metadata.Metadata {
	return m.metadata
}

func (m *messageConsumeContextBase) Created() time.Time {
	return m.created
}

func (m *messageConsumeContextBase) DeliveryTag() uint64 {
	return m.tag
}
