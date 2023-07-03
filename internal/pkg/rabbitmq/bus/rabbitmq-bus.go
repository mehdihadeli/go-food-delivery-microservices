//go:build go1.18

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
	typeMapper "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/reflection/type_mappper"
)

type RabbitmqBus interface {
	bus.Bus
	consumerConfigurations.RabbitMQConsumerConnector
}

type rabbitmqBus struct {
	messageTypeConsumers    map[reflect.Type][]consumer.Consumer
	producer                producer.Producer
	rabbitmqConfiguration   *configurations.RabbitMQConfiguration
	rabbitmqConfig          *config.RabbitmqOptions
	rabbitmqConfigBuilder   configurations.RabbitMQConfigurationBuilder
	logger                  logger.Logger
	serializer              serializer.EventSerializer
	rabbitmqConnection      types2.IConnection
	isConsumedNotifications []func(message types.IMessage)
	isProducedNotifications []func(message types.IMessage)
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

	conn, err := types2.NewRabbitMQConnection(cfg)
	if err != nil {
		return nil, err
	}

	rabbitBus := &rabbitmqBus{
		logger:                logger,
		serializer:            serializer,
		rabbitmqConfiguration: rabbitmqConfiguration,
		rabbitmqConfig:        cfg,
		rabbitmqConfigBuilder: builder,
		messageTypeConsumers:  map[reflect.Type][]consumer.Consumer{},
		rabbitmqConnection:    conn,
	}

	producersConfiguration := go2linq.ToMapMust(
		go2linq.NewEnSlice(rabbitBus.rabbitmqConfiguration.ProducersConfigurations...),
		func(source *producerConfigurations.RabbitMQProducerConfiguration) string {
			return source.ProducerMessageType.String()
		},
	)

	consumersConfiguration := go2linq.ToMapMust(
		go2linq.NewEnSlice(rabbitBus.rabbitmqConfiguration.ConsumersConfigurations...),
		func(source *consumerConfigurations.RabbitMQConsumerConfiguration) string {
			return source.ConsumerMessageType.String()
		},
	)

	for _, consumerConfiguration := range consumersConfiguration {
		mqConsumer, err := consumer2.NewRabbitMQConsumer(
			rabbitBus.rabbitmqConnection,
			consumerConfiguration,
			rabbitBus.serializer,
			rabbitBus.logger,
			// IsConsumed Notification
			func(message types.IMessage) {
				if rabbitBus.isConsumedNotifications != nil {
					for _, notification := range rabbitBus.isConsumedNotifications {
						notification(message)
					}
				}
			},
		)
		if err != nil {
			return nil, err
		}
		rabbitBus.messageTypeConsumers[consumerConfiguration.ConsumerMessageType] = append(
			rabbitBus.messageTypeConsumers[consumerConfiguration.ConsumerMessageType],
			mqConsumer,
		)
	}

	mqProducer, err := producer2.NewRabbitMQProducer(
		rabbitBus.rabbitmqConnection,
		producersConfiguration,
		rabbitBus.logger,
		rabbitBus.serializer,
		// IsProduced Notification
		func(message types.IMessage) {
			if rabbitBus.isProducedNotifications != nil {
				for _, notification := range rabbitBus.isProducedNotifications {
					notification(message)
				}
			}
		},
	)
	if err != nil {
		return nil, err
	}
	rabbitBus.producer = mqProducer

	return rabbitBus, nil
}

func (r *rabbitmqBus) IsConsumed(h func(message types.IMessage)) {
	r.isConsumedNotifications = append(r.isConsumedNotifications, h)
}

func (r *rabbitmqBus) IsProduced(h func(message types.IMessage)) {
	r.isProducedNotifications = append(r.isProducedNotifications, h)
}

// ConnectConsumer Add a new consumer to existing message type consumers. if there is no consumer, will create a new consumer for the message type
func (r *rabbitmqBus) ConnectConsumer(
	messageType types.IMessage,
	consumer consumer.Consumer,
) error {
	typeName := utils.GetMessageBaseReflectType(messageType)

	r.messageTypeConsumers[typeName] = append(r.messageTypeConsumers[typeName], consumer)

	return nil
}

// ConnectRabbitMQConsumer Add a new consumer to existing message type consumers. if there is no consumer, will create a new consumer for the message type
func (r *rabbitmqBus) ConnectRabbitMQConsumer(
	messageType types.IMessage,
	consumerBuilderFunc consumerConfigurations.RabbitMQConsumerConfigurationBuilderFuc,
) error {
	typeName := utils.GetMessageBaseReflectType(messageType)

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
		// IsConsumed Notification
		func(message types.IMessage) {
			if len(r.isConsumedNotifications) > 0 {
				for _, notification := range r.isConsumedNotifications {
					if notification != nil {
						notification(message)
					}
				}
			}
		},
	)
	if err != nil {
		return err
	}

	r.messageTypeConsumers[typeName] = append(r.messageTypeConsumers[typeName], mqConsumer)

	return nil
}

// ConnectConsumerHandler Add handler to existing consumer. creates new consumer if not exist
func (r *rabbitmqBus) ConnectConsumerHandler(
	messageType types.IMessage,
	consumerHandler consumer.ConsumerHandler,
) error {
	typeName := utils.GetMessageBaseReflectType(messageType)

	consumersForType := r.messageTypeConsumers[typeName]
	// if there is a consumer for a message type, we should add handler to existing consumers
	if consumersForType != nil {
		for _, c := range consumersForType {
			c.ConnectHandler(consumerHandler)
		}
	} else {
		// if there is no consumer for a message type, we should create new one and add handler to the consumer
		consumerBuilder := consumerConfigurations.NewRabbitMQConsumerConfigurationBuilder(messageType)
		consumerBuilder.WithHandlers(func(builder consumer.ConsumerHandlerConfigurationBuilder) {
			builder.AddHandler(consumerHandler)
		})
		consumerConfig := consumerBuilder.Build()
		mqConsumer, err := consumer2.NewRabbitMQConsumer(
			r.rabbitmqConnection,
			consumerConfig,
			r.serializer,
			r.logger,
			// IsConsumed Notification
			func(message types.IMessage) {
				if len(r.isConsumedNotifications) > 0 {
					for _, notification := range r.isConsumedNotifications {
						if notification != nil {
							notification(message)
						}
					}
				}
			},
		)
		if err != nil {
			return err
		}

		r.messageTypeConsumers[typeName] = append(r.messageTypeConsumers[typeName], mqConsumer)
	}
	return nil
}

func (r *rabbitmqBus) Start(ctx context.Context) error {
	for messageType, consumers := range r.messageTypeConsumers {
		name := typeMapper.GetTypeNameByType(messageType)
		r.logger.Info(fmt.Sprintf("consuming message type %s", name))
		for _, rabbitConsumer := range consumers {
			err := rabbitConsumer.Start(ctx)
			r.logger.Info(
				fmt.Sprintf("consumer %s, started", rabbitConsumer.GetName()),
			)
			if errors.Is(err, rabbitmqErrors.ErrDisconnected) {
				r.logger.Info(
					fmt.Sprintf(
						"consumer %s, disconnected with err: %v",
						rabbitConsumer.GetName(),
						err,
					),
				)
				// will process again with reConsume functionality
				continue
			} else if err != nil {
				r.logger.Error(
					fmt.Sprintf(
						"error in consumer %s, with err: %v",
						rabbitConsumer.GetName(),
						err,
					),
				)
				err2 := r.Stop()
				if err2 != nil {
					return errors.WrapIf(err, err2.Error())
				}
				return err
			}
		}
	}

	return nil
}

func (r *rabbitmqBus) Stop() error {
	waitGroup := sync.WaitGroup{}

	for _, consumers := range r.messageTypeConsumers {
		for _, c := range consumers {
			waitGroup.Add(1)

			go func(c consumer.Consumer) {
				defer waitGroup.Done()

				err := c.Stop()
				if err != nil {
					r.logger.Error("error in the unconsuming")
				}
			}(c)
		}
	}
	waitGroup.Wait()

	err := r.rabbitmqConnection.Close()

	return err
}

func (r *rabbitmqBus) PublishMessage(
	ctx context.Context,
	message types.IMessage,
	meta metadata.Metadata,
) error {
	if r.producer == nil {
		r.logger.Fatal("can't find a producer for publishing messages")
	}
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
