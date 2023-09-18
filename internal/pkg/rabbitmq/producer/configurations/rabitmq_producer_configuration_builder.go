package configurations

import (
	types2 "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/messaging/types"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/rabbitmq/types"
)

type RabbitMQProducerConfigurationBuilder interface {
	WithDurable(durable bool) RabbitMQProducerConfigurationBuilder
	WithAutoDeleteExchange(autoDelete bool) RabbitMQProducerConfigurationBuilder
	WithExchangeType(exchangeType types.ExchangeType) RabbitMQProducerConfigurationBuilder
	WithExchangeName(exchangeName string) RabbitMQProducerConfigurationBuilder
	WithRoutingKey(routingKey string) RabbitMQProducerConfigurationBuilder
	WithExchangeArgs(args map[string]any) RabbitMQProducerConfigurationBuilder
	WithDeliveryMode(deliveryMode uint8) RabbitMQProducerConfigurationBuilder
	WithPriority(priority uint8) RabbitMQProducerConfigurationBuilder
	WithAppId(appId string) RabbitMQProducerConfigurationBuilder
	WithExpiration(expiration string) RabbitMQProducerConfigurationBuilder
	WithReplyTo(replyTo string) RabbitMQProducerConfigurationBuilder
	WithContentEncoding(contentEncoding string) RabbitMQProducerConfigurationBuilder
	Build() *RabbitMQProducerConfiguration
}

type rabbitMQProducerConfigurationBuilder struct {
	rabbitmqProducerOptions *RabbitMQProducerConfiguration
}

func NewRabbitMQProducerConfigurationBuilder(
	messageType types2.IMessage,
) RabbitMQProducerConfigurationBuilder {
	return &rabbitMQProducerConfigurationBuilder{
		rabbitmqProducerOptions: NewDefaultRabbitMQProducerConfiguration(messageType),
	}
}

func (b *rabbitMQProducerConfigurationBuilder) WithDurable(
	durable bool,
) RabbitMQProducerConfigurationBuilder {
	b.rabbitmqProducerOptions.ExchangeOptions.Durable = durable
	return b
}

func (b *rabbitMQProducerConfigurationBuilder) WithAutoDeleteExchange(
	autoDelete bool,
) RabbitMQProducerConfigurationBuilder {
	b.rabbitmqProducerOptions.ExchangeOptions.AutoDelete = autoDelete
	return b
}

func (b *rabbitMQProducerConfigurationBuilder) WithExchangeType(
	exchangeType types.ExchangeType,
) RabbitMQProducerConfigurationBuilder {
	b.rabbitmqProducerOptions.ExchangeOptions.Type = exchangeType
	return b
}

func (b *rabbitMQProducerConfigurationBuilder) WithRoutingKey(
	routingKey string,
) RabbitMQProducerConfigurationBuilder {
	b.rabbitmqProducerOptions.RoutingKey = routingKey
	return b
}

func (b *rabbitMQProducerConfigurationBuilder) WithExchangeName(
	exchangeName string,
) RabbitMQProducerConfigurationBuilder {
	b.rabbitmqProducerOptions.ExchangeOptions.Name = exchangeName
	return b
}

func (b *rabbitMQProducerConfigurationBuilder) WithExchangeArgs(
	args map[string]any,
) RabbitMQProducerConfigurationBuilder {
	b.rabbitmqProducerOptions.ExchangeOptions.Args = args
	return b
}

func (b *rabbitMQProducerConfigurationBuilder) WithDeliveryMode(
	deliveryMode uint8,
) RabbitMQProducerConfigurationBuilder {
	b.rabbitmqProducerOptions.DeliveryMode = deliveryMode
	return b
}

func (b *rabbitMQProducerConfigurationBuilder) WithPriority(
	priority uint8,
) RabbitMQProducerConfigurationBuilder {
	b.rabbitmqProducerOptions.Priority = priority
	return b
}

func (b *rabbitMQProducerConfigurationBuilder) WithAppId(
	appId string,
) RabbitMQProducerConfigurationBuilder {
	b.rabbitmqProducerOptions.AppId = appId
	return b
}

func (b *rabbitMQProducerConfigurationBuilder) WithExpiration(
	expiration string,
) RabbitMQProducerConfigurationBuilder {
	b.rabbitmqProducerOptions.Expiration = expiration
	return b
}

func (b *rabbitMQProducerConfigurationBuilder) WithReplyTo(
	replyTo string,
) RabbitMQProducerConfigurationBuilder {
	b.rabbitmqProducerOptions.ReplyTo = replyTo
	return b
}

func (b *rabbitMQProducerConfigurationBuilder) WithContentEncoding(
	contentEncoding string,
) RabbitMQProducerConfigurationBuilder {
	b.rabbitmqProducerOptions.ContentEncoding = contentEncoding
	return b
}

func (b *rabbitMQProducerConfigurationBuilder) Build() *RabbitMQProducerConfiguration {
	return b.rabbitmqProducerOptions
}
