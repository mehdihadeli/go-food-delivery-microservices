package bus

import (
	"context"
	"emperror.dev/errors"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/core/metadata"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/core/serializer"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/bus"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/consumer"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/producer"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/types"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/configurations"
	consumer2 "github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/consumer"
	consumerConfigurations "github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/consumer/configurations"
	producer2 "github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/producer"
	producerConfigurations "github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/producer/configurations"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/rabbitmqErrors"
	types2 "github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/types"
	"github.com/solsw/go2linq/v2"
	"sync"
	"time"
)

type rabbitMQBus struct {
	consumers             []consumer.Consumer
	producer              producer.Producer
	rabbitmqConfiguration *configurations.RabbitMQConfiguration
	rabbitmqConfig        *config.RabbitMQConfig
	logger                logger.Logger
}

func AddRabbitMQBus(ctx context.Context, cfg *config.RabbitMQConfig, rabbitmqBuilderFunc configurations.RabbitMQConfigurationBuilderFuc, serializer serializer.EventSerializer, logger logger.Logger) (bus.Bus, error) {
	builder := configurations.NewRabbitMQConfigurationBuilder()
	rabbitmqBuilderFunc(builder)

	rabbitmqConfiguration := builder.Build()

	conn, err := types2.NewRabbitMQConnection(ctx, cfg)
	if err != nil {
		return nil, err
	}

	producersConfiguration := go2linq.ToMapMust(go2linq.NewEnSlice(rabbitmqConfiguration.ProducersConfigurations...), func(source *producerConfigurations.RabbitMQProducerConfiguration) string {
		return source.ProducerMessageType.String()
	})

	consumersConfiguration := go2linq.ToMapMust(go2linq.NewEnSlice(rabbitmqConfiguration.ConsumersConfigurations...), func(source *consumerConfigurations.RabbitMQConsumerConfiguration) string {
		return source.ConsumerMessageType.String()
	})

	p, err := producer2.NewRabbitMQProducer(conn, producersConfiguration, logger, serializer)
	if err != nil {
		return nil, err
	}

	var consumers []consumer.Consumer
	for _, consumerConfiguration := range consumersConfiguration {
		c, err := consumer2.NewRabbitMQConsumer(conn, consumerConfiguration, serializer, logger)
		if err != nil {
			return nil, err
		}
		consumers = append(consumers, c)
	}
	consumers[0].Start(ctx)
	time.Sleep(time.Second * 2)

	rabbitBus := &rabbitMQBus{
		logger:                logger,
		rabbitmqConfiguration: rabbitmqConfiguration,
		rabbitmqConfig:        cfg,
		producer:              p,
		consumers:             consumers,
	}

	return rabbitBus, err
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

func (r *rabbitMQBus) PublishMessage(ctx context.Context, message types.IMessage, meta metadata.Metadata) error {
	return r.producer.PublishMessage(ctx, message, meta)
}

func (r *rabbitMQBus) PublishMessageWithTopicName(ctx context.Context, message types.IMessage, meta metadata.Metadata, topicOrExchangeName string) error {
	return r.producer.PublishMessageWithTopicName(ctx, message, meta, topicOrExchangeName)
}
