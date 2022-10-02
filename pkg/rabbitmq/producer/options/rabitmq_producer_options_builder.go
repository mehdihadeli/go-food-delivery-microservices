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

func (b *RabbitMQProducerOptionsBuilder) WithExchangeName(exchangeName string) *RabbitMQProducerOptionsBuilder {
	b.rabbitmqProducerOptions.ExchangeOptions.Name = exchangeName
	return b
}

func (b *RabbitMQProducerOptionsBuilder) WithExchangeArgs(args map[string]any) *RabbitMQProducerOptionsBuilder {
	b.rabbitmqProducerOptions.ExchangeOptions.Args = args
	return b
}

func (b *RabbitMQProducerOptionsBuilder) WithDeliveryMode(deliveryMode uint8) *RabbitMQProducerOptionsBuilder {
	b.rabbitmqProducerOptions.DeliveryMode = deliveryMode
	return b
}

func (b *RabbitMQProducerOptionsBuilder) WithPriority(priority uint8) *RabbitMQProducerOptionsBuilder {
	b.rabbitmqProducerOptions.Priority = priority
	return b
}

func (b *RabbitMQProducerOptionsBuilder) WithAppId(appId string) *RabbitMQProducerOptionsBuilder {
	b.rabbitmqProducerOptions.AppId = appId
	return b
}

func (b *RabbitMQProducerOptionsBuilder) WithExpiration(expiration string) *RabbitMQProducerOptionsBuilder {
	b.rabbitmqProducerOptions.Expiration = expiration
	return b
}

func (b *RabbitMQProducerOptionsBuilder) WithReplyTo(replyTo string) *RabbitMQProducerOptionsBuilder {
	b.rabbitmqProducerOptions.ReplyTo = replyTo
	return b
}
func (b *RabbitMQProducerOptionsBuilder) WithContentEncoding(contentEncoding string) *RabbitMQProducerOptionsBuilder {
	b.rabbitmqProducerOptions.ContentEncoding = contentEncoding
	return b
}

func (b *RabbitMQProducerOptionsBuilder) Build() *RabbitMQProducerOptions {
	return b.rabbitmqProducerOptions
}
