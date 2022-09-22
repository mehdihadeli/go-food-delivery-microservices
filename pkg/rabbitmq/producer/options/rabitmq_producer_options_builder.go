package options

import "github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/types"

type RabbitMQProducerOptionsBuilder struct {
	rabbitmqProducerOptions *RabbitMQProducerOptions
}

func NewRabbitMQProducerOptionsBuilder() *RabbitMQProducerOptionsBuilder {
	return &RabbitMQProducerOptionsBuilder{rabbitmqProducerOptions: NewDefaultRabbitMQProducerOptions()}
}

func (b *RabbitMQProducerOptionsBuilder) WithDurable(durable bool) *RabbitMQProducerOptionsBuilder {
	b.rabbitmqProducerOptions.ExchangeOptions.Durable = durable
	return b
}

func (b *RabbitMQProducerOptionsBuilder) WithAutoDeleteExchange(autoDelete bool) *RabbitMQProducerOptionsBuilder {
	b.rabbitmqProducerOptions.ExchangeOptions.AutoDelete = autoDelete
	return b
}

func (b *RabbitMQProducerOptionsBuilder) WithExchangeType(exchangeType types.ExchangeType) *RabbitMQProducerOptionsBuilder {
	b.rabbitmqProducerOptions.ExchangeOptions.Type = exchangeType
	return b
}

func (b *RabbitMQProducerOptionsBuilder) WithExchangeArgs(args map[string]any) *RabbitMQProducerOptionsBuilder {
	b.rabbitmqProducerOptions.ExchangeOptions.Args = args
	return b
}

func (b *RabbitMQProducerOptionsBuilder) Build() *RabbitMQProducerOptions {
	return b.rabbitmqProducerOptions
}
