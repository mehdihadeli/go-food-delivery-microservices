package bus

import (
	"context"
	"emperror.dev/errors"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/core/metadata"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/bus"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/consumer"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/producer"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/types"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/configurations"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/rabbitmqErrors"
	"sync"
)

type rabbitMQBus struct {
	consumers             []consumer.Consumer
	producer              producer.Producer
	rabbitmqConfiguration configurations.RabbitMQConfiguration
	logger                logger.Logger
}

func NewRabbitMQBus(log logger.Logger, producer producer.Producer, consumers []consumer.Consumer) bus.Bus {
	return &rabbitMQBus{logger: log, producer: producer, consumers: consumers}
}

func (r *rabbitMQBus) Start(ctx context.Context) error {
	for _, rabbitConsumer := range r.consumers {
		err := rabbitConsumer.Consume(ctx)
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

			err := c.UnConsume(ctx)
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
