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

type RabbitMQBus interface {
	bus.Bus
	consumerConfigurations.RabbitMQConsumerConnector
}

type rabbitMQBus struct {
	consumers                map[reflect.Type]consumer.Consumer
	producer                 producer.Producer
	rabbitmqConfiguration    *configurations.RabbitMQConfiguration
	rabbitmqConfig           *config.RabbitMQConfig
	logger                   logger.Logger
	serializer               serializer.EventSerializer
	rabbitmqConnection       types2.IConnection
	messageConsumedHandlers  []func(message types.IMessage)
	messagePublishedHandlers []func(message types.IMessage)
}

func NewRabbitMQBus(
	ctx context.Context,
	cfg *config.RabbitMQConfig,
	rabbitmqBuilderFunc configurations.RabbitMQConfigurationBuilderFuc,
	serializer serializer.EventSerializer,
	logger logger.Logger,
) (RabbitMQBus, error) {
	builder := configurations.NewRabbitMQConfigurationBuilder()
	if rabbitmqBuilderFunc != nil {
		rabbitmqBuilderFunc(builder)
	}

	rabbitmqConfiguration := builder.Build()

	conn, err := types2.NewRabbitMQConnection(ctx, cfg)
	if err != nil {
		return *new(RabbitMQBus), err
	}

	producersConfiguration := go2linq.ToMapMust(
		go2linq.NewEnSlice(rabbitmqConfiguration.ProducersConfigurations...),
		func(source *producerConfigurations.RabbitMQProducerConfiguration) string {
			return source.ProducerMessageType.String()
		},
	)

	consumersConfiguration := go2linq.ToMapMust(
		go2linq.NewEnSlice(rabbitmqConfiguration.ConsumersConfigurations...),
		func(source *consumerConfigurations.RabbitMQConsumerConfiguration) string {
			return source.ConsumerMessageType.String()
		},
	)

	rabbitBus := &rabbitMQBus{
		logger:                logger,
		serializer:            serializer,
		rabbitmqConfiguration: rabbitmqConfiguration,
		rabbitmqConfig:        cfg,
		rabbitmqConnection:    conn,
	}

	var consumers = map[reflect.Type]consumer.Consumer{}
	for _, consumerConfiguration := range consumersConfiguration {
		c, err := consumer2.NewRabbitMQConsumer(
			conn,
			consumerConfiguration,
			serializer,
			logger,
			func(message types.IMessage) {
				if rabbitBus.messageConsumedHandlers != nil {
					for _, handler := range rabbitBus.messageConsumedHandlers {
						handler(message)
					}
				}
			},
		)
		if err != nil {
			return *new(RabbitMQBus), err
		}
		consumers[consumerConfiguration.ConsumerMessageType] = c
	}
	rabbitBus.consumers = consumers

	p, err := producer2.NewRabbitMQProducer(
		conn,
		producersConfiguration,
		logger,
		serializer,
		func(message types.IMessage) {
			if rabbitBus.messagePublishedHandlers != nil {
				for _, handler := range rabbitBus.messagePublishedHandlers {
					handler(message)
				}
			}
		},
	)
	if err != nil {
		return *new(RabbitMQBus), err
	}
	rabbitBus.producer = p

	return rabbitBus, err
}

func (r *rabbitMQBus) AddMessageConsumedHandler(h func(message types.IMessage)) {
	r.messageConsumedHandlers = append(r.messageConsumedHandlers, h)
}

func (r *rabbitMQBus) AddMessageProducedHandler(h func(message types.IMessage)) {
	r.messagePublishedHandlers = append(r.messagePublishedHandlers, h)
}

func (r *rabbitMQBus) ConnectConsumer(
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

func (r *rabbitMQBus) ConnectRabbitMQConsumer(
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

func (r *rabbitMQBus) ConnectConsumerHandler(
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

func (r *rabbitMQBus) Start(ctx context.Context) error {
	for _, rabbitConsumer := range r.consumers {
		err := rabbitConsumer.Start(ctx)
		if errors.Is(err, rabbitmqErrors.ErrDisconnected) {
			// will process again with reConsume functionality
			continue
		} else if err != nil {
			err2 := r.Stop(ctx)
			if err2 != nil {
				return errors.WrapIf(err, err2.Error())
			}
			return err
		}
	}

	return nil
}

func (r *rabbitMQBus) Stop(ctx context.Context) error {
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(len(r.consumers))

	for _, c := range r.consumers {
		go func(c consumer.Consumer) {
			defer waitGroup.Done()

			err := c.Stop(ctx)
			if err != nil {
				r.logger.Error("error in the unconsuming")
			}
		}(c)
	}
	waitGroup.Wait()

	return nil
}

func (r *rabbitMQBus) PublishMessage(
	ctx context.Context,
	message types.IMessage,
	meta metadata.Metadata,
) error {
	return r.producer.PublishMessage(ctx, message, meta)
}

func (r *rabbitMQBus) PublishMessageWithTopicName(
	ctx context.Context,
	message types.IMessage,
	meta metadata.Metadata,
	topicOrExchangeName string,
) error {
	return r.producer.PublishMessageWithTopicName(ctx, message, meta, topicOrExchangeName)
}
