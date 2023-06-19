//go:build.sh go1.18

package bus

import (
	"context"
	"fmt"
	"reflect"
	"sync"

	"emperror.dev/errors"
	"github.com/solsw/go2linq/v2"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/core/metadata"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/core/serializer"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/messaging/bus"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/messaging/consumer"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/messaging/producer"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/messaging/types"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/messaging/utils"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/rabbitmq/config"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/rabbitmq/configurations"
	consumer2 "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/rabbitmq/consumer"
	consumerConfigurations "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/rabbitmq/consumer/configurations"
	producer2 "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/rabbitmq/producer"
	producerConfigurations "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/rabbitmq/producer/configurations"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/rabbitmq/rabbitmqErrors"
	types2 "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/rabbitmq/types"
)

type RabbitmqBus interface {
	bus.Bus
	consumerConfigurations.RabbitMQConsumerConnector
}

type rabbitmqBus struct {
	consumers                map[reflect.Type]consumer.Consumer
	producer                 producer.Producer
	rabbitmqConfiguration    *configurations.RabbitMQConfiguration
	rabbitmqConfig           *config.RabbitmqOptions
	rabbitmqConfigBuilder    configurations.RabbitMQConfigurationBuilder
	logger                   logger.Logger
	serializer               serializer.EventSerializer
	rabbitmqConnection       types2.IConnection
	messageConsumedHandlers  []func(message types.IMessage)
	messagePublishedHandlers []func(message types.IMessage)
}

func NewRabbitmqBus(
	cfg *config.RabbitmqOptions,
	serializer serializer.EventSerializer,
	logger logger.Logger,
	rabbitmqBuilderFunc configurations.RabbitMQConfigurationBuilderFuc,
) (RabbitmqBus, error) {
	builder := configurations.NewRabbitMQConfigurationBuilder()
	if rabbitmqBuilderFunc != nil {
		rabbitmqBuilderFunc(builder)
	}

	rabbitmqConfiguration := builder.Build()

	rabbitBus := &rabbitmqBus{
		logger:                logger,
		serializer:            serializer,
		rabbitmqConfiguration: rabbitmqConfiguration,
		rabbitmqConfig:        cfg,
		rabbitmqConfigBuilder: builder,
	}

	return rabbitBus, nil
}

func (r *rabbitmqBus) AddMessageConsumedHandler(h func(message types.IMessage)) {
	r.messageConsumedHandlers = append(r.messageConsumedHandlers, h)
}

func (r *rabbitmqBus) AddMessageProducedHandler(h func(message types.IMessage)) {
	r.messagePublishedHandlers = append(r.messagePublishedHandlers, h)
}

func (r *rabbitmqBus) ConnectConsumer(
	messageType types.IMessage,
	consumer consumer.Consumer,
) error {
	c := r.consumers[utils.GetMessageBaseReflectType(messageType)]
	if c == nil {
		r.consumers[utils.GetMessageBaseReflectType(messageType)] = consumer
	} else {
		return errors.New(fmt.Sprintf("consumer %s already registerd", utils.GetMessageBaseReflectType(messageType).String()))
	}

	return nil
}

func (r *rabbitmqBus) ConnectRabbitMQConsumer(
	messageType types.IMessage,
	consumerBuilderFunc consumerConfigurations.RabbitMQConsumerConfigurationBuilderFuc,
) error {
	c := r.consumers[utils.GetMessageBaseReflectType(messageType)]
	if c == nil {
		builder := consumerConfigurations.NewRabbitMQConsumerConfigurationBuilder(messageType)
		if consumerBuilderFunc != nil {
			consumerBuilderFunc(builder)
		}
		consumerConfig := builder.Build()
		mqConsumer, err := consumer2.NewRabbitMQConsumer(
			r.rabbitmqConnection,
			consumerConfig,
			r.serializer,
			r.logger,
			func(message types.IMessage) {
				if len(r.messageConsumedHandlers) > 0 {
					for _, handler := range r.messageConsumedHandlers {
						if handler != nil {
							handler(message)
						}
					}
				}
			},
		)
		if err != nil {
			return err
		}
		r.consumers[utils.GetMessageBaseReflectType(messageType)] = mqConsumer
	} else {
		return errors.New(fmt.Sprintf("consumer %s already registerd", utils.GetMessageBaseReflectType(messageType).String()))
	}

	return nil
}

func (r *rabbitmqBus) ConnectConsumerHandler(
	messageType types.IMessage,
	consumerHandler consumer.ConsumerHandler,
) error {
	c := r.consumers[utils.GetMessageBaseReflectType(messageType)]
	if c != nil {
		c.ConnectHandler(consumerHandler)
	} else {
		consumerBuilder := consumerConfigurations.NewRabbitMQConsumerConfigurationBuilder(messageType)
		consumerBuilder.WithHandlers(func(builder consumer.ConsumerHandlerConfigurationBuilder) {
			builder.AddHandler(consumerHandler)
		})
		consumerConfig := consumerBuilder.Build()
		mqConsumer, err := consumer2.NewRabbitMQConsumer(r.rabbitmqConnection, consumerConfig, r.serializer, r.logger, func(message types.IMessage) {
			if len(r.messageConsumedHandlers) > 0 {
				for _, handler := range r.messageConsumedHandlers {
					if handler != nil {
						handler(message)
					}
				}
			}
		})
		if err != nil {
			return err
		}
		r.consumers[utils.GetMessageBaseReflectType(messageType)] = mqConsumer
	}
	return nil
}

func (r *rabbitmqBus) Start(ctx context.Context) error {
	conn, err := types2.NewRabbitMQConnection(ctx, r.rabbitmqConfig)
	if err != nil {
		return err
	}
	r.rabbitmqConnection = conn

	producersConfiguration := go2linq.ToMapMust(
		go2linq.NewEnSlice(r.rabbitmqConfiguration.ProducersConfigurations...),
		func(source *producerConfigurations.RabbitMQProducerConfiguration) string {
			return source.ProducerMessageType.String()
		},
	)

	consumersConfiguration := go2linq.ToMapMust(
		go2linq.NewEnSlice(r.rabbitmqConfiguration.ConsumersConfigurations...),
		func(source *consumerConfigurations.RabbitMQConsumerConfiguration) string {
			return source.ConsumerMessageType.String()
		},
	)

	consumers := map[reflect.Type]consumer.Consumer{}
	for _, consumerConfiguration := range consumersConfiguration {
		c, err := consumer2.NewRabbitMQConsumer(
			conn,
			consumerConfiguration,
			r.serializer,
			r.logger,
			func(message types.IMessage) {
				if r.messageConsumedHandlers != nil {
					for _, handler := range r.messageConsumedHandlers {
						handler(message)
					}
				}
			},
		)
		if err != nil {
			return err
		}
		consumers[consumerConfiguration.ConsumerMessageType] = c
	}
	r.consumers = consumers

	p, err := producer2.NewRabbitMQProducer(
		conn,
		producersConfiguration,
		r.logger,
		r.serializer,
		func(message types.IMessage) {
			if r.messagePublishedHandlers != nil {
				for _, handler := range r.messagePublishedHandlers {
					handler(message)
				}
			}
		},
	)
	if err != nil {
		return err
	}
	r.producer = p

	for _, rabbitConsumer := range r.consumers {
		err := rabbitConsumer.Start(ctx)
		if errors.Is(err, rabbitmqErrors.ErrDisconnected) {
			// will process again with reConsume functionality
			continue
		} else if err != nil {
			err2 := r.Stop()
			if err2 != nil {
				return errors.WrapIf(err, err2.Error())
			}
			return err
		}
	}

	return nil
}

func (r *rabbitmqBus) Stop() error {
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(len(r.consumers))

	for _, c := range r.consumers {
		go func(c consumer.Consumer) {
			defer waitGroup.Done()

			err := c.Stop()
			if err != nil {
				r.logger.Error("error in the unconsuming")
			}
		}(c)
	}
	waitGroup.Wait()

	return nil
}

func (r *rabbitmqBus) PublishMessage(
	ctx context.Context,
	message types.IMessage,
	meta metadata.Metadata,
) error {
	return r.producer.PublishMessage(ctx, message, meta)
}

func (r *rabbitmqBus) PublishMessageWithTopicName(
	ctx context.Context,
	message types.IMessage,
	meta metadata.Metadata,
	topicOrExchangeName string,
) error {
	return r.producer.PublishMessageWithTopicName(ctx, message, meta, topicOrExchangeName)
}
