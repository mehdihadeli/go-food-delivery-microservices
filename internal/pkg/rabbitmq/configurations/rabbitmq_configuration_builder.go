//go:build go1.18

package configurations

import (
	"github.com/solsw/go2linq/v2"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/messaging/types"
	consumerConfigurations "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/rabbitmq/consumer/configurations"
	producerConfigurations "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/rabbitmq/producer/configurations"
)

type RabbitMQConfigurationBuilder interface {
	AddProducer(
		producerMessageType types.IMessage,
		producerBuilderFunc producerConfigurations.RabbitMQProducerConfigurationBuilderFuc,
	) RabbitMQConfigurationBuilder
	AddConsumer(
		consumerMessageType types.IMessage,
		consumerBuilderFunc consumerConfigurations.RabbitMQConsumerConfigurationBuilderFuc,
	) RabbitMQConfigurationBuilder
	Build() *RabbitMQConfiguration
}

type rabbitMQConfigurationBuilder struct {
	rabbitMQConfiguration *RabbitMQConfiguration
	consumerBuilders      []consumerConfigurations.RabbitMQConsumerConfigurationBuilder
	producerBuilders      []producerConfigurations.RabbitMQProducerConfigurationBuilder
}

func NewRabbitMQConfigurationBuilder() RabbitMQConfigurationBuilder {
	return &rabbitMQConfigurationBuilder{
		rabbitMQConfiguration: &RabbitMQConfiguration{},
	}
}

func (r *rabbitMQConfigurationBuilder) AddProducer(
	producerMessageType types.IMessage,
	producerBuilderFunc producerConfigurations.RabbitMQProducerConfigurationBuilderFuc,
) RabbitMQConfigurationBuilder {
	builder := producerConfigurations.NewRabbitMQProducerConfigurationBuilder(producerMessageType)
	if producerBuilderFunc != nil {
		producerBuilderFunc(builder)
	}
	r.producerBuilders = append(r.producerBuilders, builder)

	return r
}

func (r *rabbitMQConfigurationBuilder) AddConsumer(
	consumerMessageType types.IMessage,
	consumerBuilderFunc consumerConfigurations.RabbitMQConsumerConfigurationBuilderFuc,
) RabbitMQConfigurationBuilder {
	builder := consumerConfigurations.NewRabbitMQConsumerConfigurationBuilder(consumerMessageType)
	if consumerBuilderFunc != nil {
		consumerBuilderFunc(builder)
	}
	r.consumerBuilders = append(r.consumerBuilders, builder)

	return r
}

func (r *rabbitMQConfigurationBuilder) Build() *RabbitMQConfiguration {
	consumersConfig := go2linq.ToSliceMust(
		go2linq.SelectMust(
			go2linq.NewEnSlice(r.consumerBuilders...),
			func(source consumerConfigurations.RabbitMQConsumerConfigurationBuilder) *consumerConfigurations.RabbitMQConsumerConfiguration {
				return source.Build()
			},
		))

	producersConfig := go2linq.ToSliceMust(
		go2linq.SelectMust(
			go2linq.NewEnSlice(r.producerBuilders...),
			func(source producerConfigurations.RabbitMQProducerConfigurationBuilder) *producerConfigurations.RabbitMQProducerConfiguration {
				return source.Build()
			},
		))

	r.rabbitMQConfiguration.ConsumersConfigurations = consumersConfig
	r.rabbitMQConfiguration.ProducersConfigurations = producersConfig

	return r.rabbitMQConfiguration
}
