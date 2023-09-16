package configurations

import (
	messageConsumer "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/messaging/consumer"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/messaging/pipeline"
	types2 "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/messaging/types"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/rabbitmq/types"
)

type RabbitMQConsumerConfigurationBuilder interface {
	WithHandlers(
		consumerBuilderFunc messageConsumer.ConsumerHandlerConfigurationBuilderFunc,
	) RabbitMQConsumerConfigurationBuilder
	WIthPipelines(
		pipelineBuilderFunc pipeline.ConsumerPipelineConfigurationBuilderFunc,
	) RabbitMQConsumerConfigurationBuilder
	WithExitOnError(exitOnError bool) RabbitMQConsumerConfigurationBuilder
	WithAutoAck(ack bool) RabbitMQConsumerConfigurationBuilder
	WithNoLocal(noLocal bool) RabbitMQConsumerConfigurationBuilder
	WithNoWait(noWait bool) RabbitMQConsumerConfigurationBuilder
	WithConcurrencyLimit(limit int) RabbitMQConsumerConfigurationBuilder
	WithPrefetchCount(count int) RabbitMQConsumerConfigurationBuilder
	WithConsumerId(consumerId string) RabbitMQConsumerConfigurationBuilder
	WithQueueName(queueName string) RabbitMQConsumerConfigurationBuilder
	WithDurable(durable bool) RabbitMQConsumerConfigurationBuilder
	WithAutoDeleteQueue(autoDelete bool) RabbitMQConsumerConfigurationBuilder
	WithExclusiveQueue(exclusive bool) RabbitMQConsumerConfigurationBuilder
	WithQueueArgs(args map[string]any) RabbitMQConsumerConfigurationBuilder
	WithExchangeName(exchangeName string) RabbitMQConsumerConfigurationBuilder
	WithAutoDeleteExchange(autoDelete bool) RabbitMQConsumerConfigurationBuilder
	WithExchangeType(exchangeType types.ExchangeType) RabbitMQConsumerConfigurationBuilder
	WithExchangeArgs(args map[string]any) RabbitMQConsumerConfigurationBuilder
	WithRoutingKey(routingKey string) RabbitMQConsumerConfigurationBuilder
	WithBindingArgs(args map[string]any) RabbitMQConsumerConfigurationBuilder
	WithName(name string) RabbitMQConsumerConfigurationBuilder
	Build() *RabbitMQConsumerConfiguration
}

type rabbitMQConsumerConfigurationBuilder struct {
	rabbitmqConsumerConfigurations *RabbitMQConsumerConfiguration
	pipelinesBuilder               pipeline.ConsumerPipelineConfigurationBuilder
	handlersBuilder                messageConsumer.ConsumerHandlerConfigurationBuilder
}

func NewRabbitMQConsumerConfigurationBuilder(
	messageType types2.IMessage,
) RabbitMQConsumerConfigurationBuilder {
	return &rabbitMQConsumerConfigurationBuilder{
		rabbitmqConsumerConfigurations: NewDefaultRabbitMQConsumerConfiguration(messageType),
	}
}

func (b *rabbitMQConsumerConfigurationBuilder) WIthPipelines(
	pipelineBuilderFunc pipeline.ConsumerPipelineConfigurationBuilderFunc,
) RabbitMQConsumerConfigurationBuilder {
	builder := pipeline.NewConsumerPipelineConfigurationBuilder()
	if pipelineBuilderFunc != nil {
		pipelineBuilderFunc(builder)
	}
	b.pipelinesBuilder = builder

	return b
}

func (b *rabbitMQConsumerConfigurationBuilder) WithHandlers(
	consumerBuilderFunc messageConsumer.ConsumerHandlerConfigurationBuilderFunc,
) RabbitMQConsumerConfigurationBuilder {
	builder := messageConsumer.NewConsumerHandlersConfigurationBuilder()
	if consumerBuilderFunc != nil {
		consumerBuilderFunc(builder)
	}
	b.handlersBuilder = builder

	return b
}

func (b *rabbitMQConsumerConfigurationBuilder) WithExitOnError(
	exitOnError bool,
) RabbitMQConsumerConfigurationBuilder {
	b.rabbitmqConsumerConfigurations.ExitOnError = exitOnError
	return b
}

func (b *rabbitMQConsumerConfigurationBuilder) WithName(
	name string,
) RabbitMQConsumerConfigurationBuilder {
	b.rabbitmqConsumerConfigurations.Name = name
	return b
}

func (b *rabbitMQConsumerConfigurationBuilder) WithAutoAck(
	ack bool,
) RabbitMQConsumerConfigurationBuilder {
	b.rabbitmqConsumerConfigurations.AutoAck = ack
	return b
}

func (b *rabbitMQConsumerConfigurationBuilder) WithNoLocal(
	noLocal bool,
) RabbitMQConsumerConfigurationBuilder {
	b.rabbitmqConsumerConfigurations.NoLocal = noLocal
	return b
}

func (b *rabbitMQConsumerConfigurationBuilder) WithNoWait(
	noWait bool,
) RabbitMQConsumerConfigurationBuilder {
	b.rabbitmqConsumerConfigurations.NoWait = noWait
	return b
}

func (b *rabbitMQConsumerConfigurationBuilder) WithConcurrencyLimit(
	limit int,
) RabbitMQConsumerConfigurationBuilder {
	b.rabbitmqConsumerConfigurations.ConcurrencyLimit = limit
	return b
}

func (b *rabbitMQConsumerConfigurationBuilder) WithPrefetchCount(
	count int,
) RabbitMQConsumerConfigurationBuilder {
	b.rabbitmqConsumerConfigurations.PrefetchCount = count
	return b
}

func (b *rabbitMQConsumerConfigurationBuilder) WithConsumerId(
	consumerId string,
) RabbitMQConsumerConfigurationBuilder {
	b.rabbitmqConsumerConfigurations.ConsumerId = consumerId
	return b
}

func (b *rabbitMQConsumerConfigurationBuilder) WithQueueName(
	queueName string,
) RabbitMQConsumerConfigurationBuilder {
	b.rabbitmqConsumerConfigurations.QueueOptions.Name = queueName
	return b
}

func (b *rabbitMQConsumerConfigurationBuilder) WithDurable(
	durable bool,
) RabbitMQConsumerConfigurationBuilder {
	b.rabbitmqConsumerConfigurations.ExchangeOptions.Durable = durable
	b.rabbitmqConsumerConfigurations.QueueOptions.Durable = durable
	return b
}

func (b *rabbitMQConsumerConfigurationBuilder) WithAutoDeleteQueue(
	autoDelete bool,
) RabbitMQConsumerConfigurationBuilder {
	b.rabbitmqConsumerConfigurations.QueueOptions.AutoDelete = autoDelete
	return b
}

func (b *rabbitMQConsumerConfigurationBuilder) WithExclusiveQueue(
	exclusive bool,
) RabbitMQConsumerConfigurationBuilder {
	b.rabbitmqConsumerConfigurations.QueueOptions.Exclusive = exclusive
	return b
}

func (b *rabbitMQConsumerConfigurationBuilder) WithQueueArgs(
	args map[string]any,
) RabbitMQConsumerConfigurationBuilder {
	b.rabbitmqConsumerConfigurations.QueueOptions.Args = args
	return b
}

func (b *rabbitMQConsumerConfigurationBuilder) WithExchangeName(
	exchangeName string,
) RabbitMQConsumerConfigurationBuilder {
	b.rabbitmqConsumerConfigurations.ExchangeOptions.Name = exchangeName
	return b
}

func (b *rabbitMQConsumerConfigurationBuilder) WithAutoDeleteExchange(
	autoDelete bool,
) RabbitMQConsumerConfigurationBuilder {
	b.rabbitmqConsumerConfigurations.ExchangeOptions.AutoDelete = autoDelete
	return b
}

func (b *rabbitMQConsumerConfigurationBuilder) WithExchangeType(
	exchangeType types.ExchangeType,
) RabbitMQConsumerConfigurationBuilder {
	b.rabbitmqConsumerConfigurations.ExchangeOptions.Type = exchangeType
	return b
}

func (b *rabbitMQConsumerConfigurationBuilder) WithExchangeArgs(
	args map[string]any,
) RabbitMQConsumerConfigurationBuilder {
	b.rabbitmqConsumerConfigurations.ExchangeOptions.Args = args
	return b
}

func (b *rabbitMQConsumerConfigurationBuilder) WithRoutingKey(
	routingKey string,
) RabbitMQConsumerConfigurationBuilder {
	b.rabbitmqConsumerConfigurations.BindingOptions.RoutingKey = routingKey
	return b
}

func (b *rabbitMQConsumerConfigurationBuilder) WithBindingArgs(
	args map[string]any,
) RabbitMQConsumerConfigurationBuilder {
	b.rabbitmqConsumerConfigurations.BindingOptions.Args = args
	return b
}

func (b *rabbitMQConsumerConfigurationBuilder) Build() *RabbitMQConsumerConfiguration {
	if b.pipelinesBuilder != nil {
		b.rabbitmqConsumerConfigurations.Pipelines = b.pipelinesBuilder.Build().Pipelines
	}
	if b.handlersBuilder != nil {
		b.rabbitmqConsumerConfigurations.Handlers = b.handlersBuilder.Build().Handlers
	}

	return b.rabbitmqConsumerConfigurations
}
