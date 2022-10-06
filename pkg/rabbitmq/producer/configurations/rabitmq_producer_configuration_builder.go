package configurations

import (
	types2 "github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/types"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/utils"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/types"
)

type RabbitMQProducerConfigurationBuilder interface {
	SetProducerMessageType(messageType types2.IMessage) RabbitMQProducerConfigurationBuilder
	WithDurable(durable bool) RabbitMQProducerConfigurationBuilder
	WithAutoDeleteExchange(autoDelete bool) RabbitMQProducerConfigurationBuilder
	WithExchangeType(exchangeType types.ExchangeType) RabbitMQProducerConfigurationBuilder
	WithExchangeName(exchangeName string) RabbitMQProducerConfigurationBuilder
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

func NewRabbitMQProducerConfigurationBuilder() RabbitMQProducerConfigurationBuilder {
	return &rabbitMQProducerConfigurationBuilder{rabbitmqProducerOptions: NewDefaultRabbitMQProducerConfiguration()}
}

func (b *rabbitMQProducerConfigurationBuilder) SetProducerMessageType(messageType types2.IMessage) RabbitMQProducerConfigurationBuilder {
	b.rabbitmqProducerOptions.ProducerMessageType = utils.GetMessageBaseReflectType(messageType)
	return b
}

func (b *rabbitMQProducerConfigurationBuilder) WithDurable(durable bool) RabbitMQProducerConfigurationBuilder {
	b.rabbitmqProducerOptions.ExchangeOptions.Durable = durable
	return b
}

func (b *rabbitMQProducerConfigurationBuilder) WithAutoDeleteExchange(autoDelete bool) RabbitMQProducerConfigurationBuilder {
	b.rabbitmqProducerOptions.ExchangeOptions.AutoDelete = autoDelete
	return b
}

func (b *rabbitMQProducerConfigurationBuilder) WithExchangeType(exchangeType types.ExchangeType) RabbitMQProducerConfigurationBuilder {
	b.rabbitmqProducerOptions.ExchangeOptions.Type = exchangeType
	return b
}

func (b *rabbitMQProducerConfigurationBuilder) WithExchangeName(exchangeName string) RabbitMQProducerConfigurationBuilder {
	b.rabbitmqProducerOptions.ExchangeOptions.Name = exchangeName
	return b
}

func (b *rabbitMQProducerConfigurationBuilder) WithExchangeArgs(args map[string]any) RabbitMQProducerConfigurationBuilder {
	b.rabbitmqProducerOptions.ExchangeOptions.Args = args
	return b
}

func (b *rabbitMQProducerConfigurationBuilder) WithDeliveryMode(deliveryMode uint8) RabbitMQProducerConfigurationBuilder {
	b.rabbitmqProducerOptions.DeliveryMode = deliveryMode
	return b
}

func (b *rabbitMQProducerConfigurationBuilder) WithPriority(priority uint8) RabbitMQProducerConfigurationBuilder {
	b.rabbitmqProducerOptions.Priority = priority
	return b
}

func (b *rabbitMQProducerConfigurationBuilder) WithAppId(appId string) RabbitMQProducerConfigurationBuilder {
	b.rabbitmqProducerOptions.AppId = appId
	return b
}

func (b *rabbitMQProducerConfigurationBuilder) WithExpiration(expiration string) RabbitMQProducerConfigurationBuilder {
	b.rabbitmqProducerOptions.Expiration = expiration
	return b
}

func (b *rabbitMQProducerConfigurationBuilder) WithReplyTo(replyTo string) RabbitMQProducerConfigurationBuilder {
	b.rabbitmqProducerOptions.ReplyTo = replyTo
	return b
}
func (b *rabbitMQProducerConfigurationBuilder) WithContentEncoding(contentEncoding string) RabbitMQProducerConfigurationBuilder {
	b.rabbitmqProducerOptions.ContentEncoding = contentEncoding
	return b
}

func (b *rabbitMQProducerConfigurationBuilder) Build() *RabbitMQProducerConfiguration {
	return b.rabbitmqProducerOptions
}
