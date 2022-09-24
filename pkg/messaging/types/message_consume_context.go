package types

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/core"
	"time"
)

type IMessageConsumeContext[T IMessage] interface {
	MessageId() string
	CorrelationId() string
	MessageType() string
	Created() time.Time
	ContentType() string
	Tag() uint64
	Metadata() core.Metadata
	Message() T
}

type messageConsumeContext[T IMessage] struct {
	message       T
	metadata      core.Metadata
	contentType   string
	messageType   string
	messageId     string
	created       time.Time
	tag           uint64
	correlationId string
}

func NewMessageConsumeContext[T IMessage](message T, meta core.Metadata, contentType string, messageType string, created time.Time, deliveryTag uint64, messageId string, correlationId string) IMessageConsumeContext[T] {
	return &messageConsumeContext[T]{
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

func (m messageConsumeContext[T]) Message() T {
	return m.message
}

func (m messageConsumeContext[T]) MessageId() string {
	return m.messageId
}

func (m messageConsumeContext[T]) CorrelationId() string {
	return m.correlationId
}

func (m messageConsumeContext[T]) MessageType() string {
	return m.messageType
}

func (m messageConsumeContext[T]) ContentType() string {
	return m.contentType
}

func (m messageConsumeContext[T]) Metadata() core.Metadata {
	return m.metadata
}

func (m messageConsumeContext[T]) Created() time.Time {
	return m.created
}

func (m messageConsumeContext[T]) Tag() uint64 {
	return m.tag
}
