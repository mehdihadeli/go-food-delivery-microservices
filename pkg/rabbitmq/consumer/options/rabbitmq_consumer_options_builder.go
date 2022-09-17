package options

import (
	types2 "github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/types"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/types"
)

type RabbitMQConsumerOptionsBuilder[T types2.IMessage] struct {
	rabbitmqConsumerOptions *RabbitMQConsumerOptions
}

func NewRabbitMQConsumerOptionsBuilder[T types2.IMessage]() *RabbitMQConsumerOptionsBuilder[T] {
	return &RabbitMQConsumerOptionsBuilder[T]{rabbitmqConsumerOptions: NewDefaultRabbitMQConsumerOptions[T]()}
}

func (b *RabbitMQConsumerOptionsBuilder[T]) WithExitOnError(exitOnError bool) *RabbitMQConsumerOptionsBuilder[T] {
	b.rabbitmqConsumerOptions.ExitOnError = exitOnError
	return b
}

func (b *RabbitMQConsumerOptionsBuilder[T]) WithAutoAck(ack bool) *RabbitMQConsumerOptionsBuilder[T] {
	b.rabbitmqConsumerOptions.AutoAck = ack
	return b
}

func (b *RabbitMQConsumerOptionsBuilder[T]) WithNoLocal(noLocal bool) *RabbitMQConsumerOptionsBuilder[T] {
	b.rabbitmqConsumerOptions.NoLocal = noLocal
	return b
}

func (b *RabbitMQConsumerOptionsBuilder[T]) WithNoWait(noWait bool) *RabbitMQConsumerOptionsBuilder[T] {
	b.rabbitmqConsumerOptions.NoWait = noWait
	return b
}

func (b *RabbitMQConsumerOptionsBuilder[T]) WithConcurrencyLimit(limit int) *RabbitMQConsumerOptionsBuilder[T] {
	b.rabbitmqConsumerOptions.ConcurrencyLimit = limit
	return b
}

func (b *RabbitMQConsumerOptionsBuilder[T]) WithPrefetchCount(count int) *RabbitMQConsumerOptionsBuilder[T] {
	b.rabbitmqConsumerOptions.PrefetchCount = count
	return b
}

func (b *RabbitMQConsumerOptionsBuilder[T]) WithConsumerId(consumerId string) *RabbitMQConsumerOptionsBuilder[T] {
	b.rabbitmqConsumerOptions.ConsumerId = consumerId
	return b
}

func (b *RabbitMQConsumerOptionsBuilder[T]) WithQueueName(queueName string) *RabbitMQConsumerOptionsBuilder[T] {
	b.rabbitmqConsumerOptions.QueueOptions.Name = queueName
	return b
}

func (b *RabbitMQConsumerOptionsBuilder[T]) WithDurable(durable bool) *RabbitMQConsumerOptionsBuilder[T] {
	b.rabbitmqConsumerOptions.ExchangeOptions.Durable = durable
	b.rabbitmqConsumerOptions.QueueOptions.Durable = durable
	return b
}

func (b *RabbitMQConsumerOptionsBuilder[T]) WithAutoDeleteQueue(autoDelete bool) *RabbitMQConsumerOptionsBuilder[T] {
	b.rabbitmqConsumerOptions.QueueOptions.AutoDelete = autoDelete
	return b
}

func (b *RabbitMQConsumerOptionsBuilder[T]) WithExclusiveQueue(exclusive bool) *RabbitMQConsumerOptionsBuilder[T] {
	b.rabbitmqConsumerOptions.QueueOptions.Exclusive = exclusive
	return b
}

func (b *RabbitMQConsumerOptionsBuilder[T]) WithQueueArgs(args map[string]any) *RabbitMQConsumerOptionsBuilder[T] {
	b.rabbitmqConsumerOptions.QueueOptions.Args = args
	return b
}

func (b *RabbitMQConsumerOptionsBuilder[T]) WithExchangeName(exchangeName string) *RabbitMQConsumerOptionsBuilder[T] {
	b.rabbitmqConsumerOptions.ExchangeOptions.Name = exchangeName
	return b
}

func (b *RabbitMQConsumerOptionsBuilder[T]) WithAutoDeleteExchange(autoDelete bool) *RabbitMQConsumerOptionsBuilder[T] {
	b.rabbitmqConsumerOptions.ExchangeOptions.AutoDelete = autoDelete
	return b
}

func (b *RabbitMQConsumerOptionsBuilder[T]) WithExchangeType(exchangeType types.ExchangeType) *RabbitMQConsumerOptionsBuilder[T] {
	b.rabbitmqConsumerOptions.ExchangeOptions.Type = exchangeType
	return b
}

func (b *RabbitMQConsumerOptionsBuilder[T]) WithExchangeArgs(args map[string]any) *RabbitMQConsumerOptionsBuilder[T] {
	b.rabbitmqConsumerOptions.ExchangeOptions.Args = args
	return b
}

func (b *RabbitMQConsumerOptionsBuilder[T]) WithRoutingKey(routingKey string) *RabbitMQConsumerOptionsBuilder[T] {
	b.rabbitmqConsumerOptions.BindingOptions.RoutingKey = routingKey
	return b
}

func (b *RabbitMQConsumerOptionsBuilder[T]) WithBindingArgs(args map[string]any) *RabbitMQConsumerOptionsBuilder[T] {
	b.rabbitmqConsumerOptions.BindingOptions.Args = args
	return b
}

func (b *RabbitMQConsumerOptionsBuilder[T]) Build() *RabbitMQConsumerOptions {
	return b.rabbitmqConsumerOptions
}
