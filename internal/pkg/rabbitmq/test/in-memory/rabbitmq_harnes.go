package in_memory

import (
	"context"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/core/metadata"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/messaging/consumer"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/messaging/types"
)

type RabbitmqInMemoryHarnesses struct {
	publishedMessage []types.IMessage
	consumedMessage  []types.IMessage
	consumerHandlers map[types.IMessage][]consumer.ConsumerHandler
}

func NewRabbitmqInMemoryHarnesses() *RabbitmqInMemoryHarnesses {
	return &RabbitmqInMemoryHarnesses{}
}

func (r *RabbitmqInMemoryHarnesses) PublishMessage(ctx context.Context, message types.IMessage, meta metadata.Metadata) error {
	r.publishedMessage = append(r.publishedMessage, message)
	return nil
}

func (r *RabbitmqInMemoryHarnesses) PublishMessageWithTopicName(ctx context.Context, message types.IMessage, meta metadata.Metadata, topicOrExchangeName string) error {
	r.publishedMessage = append(r.publishedMessage, message)
	return nil
}

func (r *RabbitmqInMemoryHarnesses) AddMessageProducedHandler(f func(message types.IMessage)) {
}

func (r *RabbitmqInMemoryHarnesses) AddMessageConsumedHandler(f func(message types.IMessage)) {
}

func (r *RabbitmqInMemoryHarnesses) Start(ctx context.Context) error {
	return nil
}

func (r *RabbitmqInMemoryHarnesses) Stop(ctx context.Context) error {
	return nil
}

func (r *RabbitmqInMemoryHarnesses) ConnectConsumerHandler(messageType types.IMessage, consumerHandler consumer.ConsumerHandler) error {
	r.consumerHandlers[messageType] = append(r.consumerHandlers[messageType], consumerHandler)
	return nil
}

func (r *RabbitmqInMemoryHarnesses) ConnectConsumer(messageType types.IMessage, consumer consumer.Consumer) error {
	return nil
}

func (r *RabbitmqInMemoryHarnesses) PublishedMessages() []types.IMessage {
	return r.publishedMessage
}

func (r *RabbitmqInMemoryHarnesses) ConsumedMessages() []types.IMessage {
	return r.consumedMessage
}
