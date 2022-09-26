package options

import (
	types2 "github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/types"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/types"
)

type RabbitMQConsumerOptionsBuilder struct {
	rabbitmqConsumerOptions *RabbitMQConsumerOptions
}

func NewRabbitMQConsumerOptionsBuilderT[T types2.IMessage]() *RabbitMQConsumerOptionsBuilder {
	return &RabbitMQConsumerOptionsBuilder{rabbitmqConsumerOptions: NewDefaultRabbitMQConsumerOptionsT[T]()}
}

func NewRabbitMQConsumerOptionsBuilder() *RabbitMQConsumerOptionsBuilder {
	return &RabbitMQConsumerOptionsBuilder{rabbitmqConsumerOptions: NewDefaultRabbitMQConsumerOptions()}
}

func (b *RabbitMQConsumerOptionsBuilder) WithExitOnError(exitOnError bool) *RabbitMQConsumerOptionsBuilder {
	b.rabbitmqConsumerOptions.ExitOnError = exitOnError
	return b
}

func (b *RabbitMQConsumerOptionsBuilder) WithAutoAck(ack bool) *RabbitMQConsumerOptionsBuilder {
	b.rabbitmqConsumerOptions.AutoAck = ack
	return b
}

func (b *RabbitMQConsumerOptionsBuilder) WithNoLocal(noLocal bool) *RabbitMQConsumerOptionsBuilder {
	b.rabbitmqConsumerOptions.NoLocal = noLocal
	return b
}

func (b *RabbitMQConsumerOptionsBuilder) WithNoWait(noWait bool) *RabbitMQConsumerOptionsBuilder {
	b.rabbitmqConsumerOptions.NoWait = noWait
	return b
}

func (b *RabbitMQConsumerOptionsBuilder) WithConcurrencyLimit(limit int) *RabbitMQConsumerOptionsBuilder {
	b.rabbitmqConsumerOptions.ConcurrencyLimit = limit
	return b
}

func (b *RabbitMQConsumerOptionsBuilder) WithPrefetchCount(count int) *RabbitMQConsumerOptionsBuilder {
	b.rabbitmqConsumerOptions.PrefetchCount = count
	return b
}

func (b *RabbitMQConsumerOptionsBuilder) WithConsumerId(consumerId string) *RabbitMQConsumerOptionsBuilder {
	b.rabbitmqConsumerOptions.ConsumerId = consumerId
	return b
}

func (b *RabbitMQConsumerOptionsBuilder) WithQueueName(queueName string) *RabbitMQConsumerOptionsBuilder {
	b.rabbitmqConsumerOptions.QueueOptions.Name = queueName
	return b
}

func (b *RabbitMQConsumerOptionsBuilder) WithDurable(durable bool) *RabbitMQConsumerOptionsBuilder {
	b.rabbitmqConsumerOptions.ExchangeOptions.Durable = durable
	b.rabbitmqConsumerOptions.QueueOptions.Durable = durable
	return b
}

func (b *RabbitMQConsumerOptionsBuilder) WithAutoDeleteQueue(autoDelete bool) *RabbitMQConsumerOptionsBuilder {
	b.rabbitmqConsumerOptions.QueueOptions.AutoDelete = autoDelete
	return b
}

func (b *RabbitMQConsumerOptionsBuilder) WithExclusiveQueue(exclusive bool) *RabbitMQConsumerOptionsBuilder {
	b.rabbitmqConsumerOptions.QueueOptions.Exclusive = exclusive
	return b
}

func (b *RabbitMQConsumerOptionsBuilder) WithQueueArgs(args map[string]any) *RabbitMQConsumerOptionsBuilder {
	b.rabbitmqConsumerOptions.QueueOptions.Args = args
	return b
}

func (b *RabbitMQConsumerOptionsBuilder) WithExchangeName(exchangeName string) *RabbitMQConsumerOptionsBuilder {
	b.rabbitmqConsumerOptions.ExchangeOptions.Name = exchangeName
	return b
}

func (b *RabbitMQConsumerOptionsBuilder) WithAutoDeleteExchange(autoDelete bool) *RabbitMQConsumerOptionsBuilder {
	b.rabbitmqConsumerOptions.ExchangeOptions.AutoDelete = autoDelete
	return b
}

func (b *RabbitMQConsumerOptionsBuilder) WithExchangeType(exchangeType types.ExchangeType) *RabbitMQConsumerOptionsBuilder {
	b.rabbitmqConsumerOptions.ExchangeOptions.Type = exchangeType
	return b
}

func (b *RabbitMQConsumerOptionsBuilder) WithExchangeArgs(args map[string]any) *RabbitMQConsumerOptionsBuilder {
	b.rabbitmqConsumerOptions.ExchangeOptions.Args = args
	return b
}

func (b *RabbitMQConsumerOptionsBuilder) WithRoutingKey(routingKey string) *RabbitMQConsumerOptionsBuilder {
	b.rabbitmqConsumerOptions.BindingOptions.RoutingKey = routingKey
	return b
}

func (b *RabbitMQConsumerOptionsBuilder) WithBindingArgs(args map[string]any) *RabbitMQConsumerOptionsBuilder {
	b.rabbitmqConsumerOptions.BindingOptions.Args = args
	return b
}

func (b *RabbitMQConsumerOptionsBuilder) Build() *RabbitMQConsumerOptions {
	return b.rabbitmqConsumerOptions
}
