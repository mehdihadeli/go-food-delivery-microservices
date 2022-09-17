package consumer

import (
	"context"
	"emperror.dev/errors"
	"github.com/avast/retry-go"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/core"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/core/serializer"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/consumer"
	types2 "github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/types"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/consumer/options"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/types"
	"github.com/rabbitmq/amqp091-go"
	"time"
)

const (
	retryAttempts = 3
	retryDelay    = 300 * time.Millisecond
)

var (
	retryOptions = []retry.Option{retry.Attempts(retryAttempts), retry.Delay(retryDelay), retry.DelayType(retry.BackOffDelay)}
)

type RabbitMQConsumer[T types2.IMessage] struct {
	rabbitmqConsumerOptions *options.RabbitMQConsumerOptions
	connection              types.IConnection
	handler                 consumer.ConsumerHandler[T]
	channel                 *amqp091.Channel
	deliveryRoutines        chan struct{}
	eventSerializer         serializer.EventSerializer
	logger                  logger.Logger
}

func NewRabbitMQConsumer[T types2.IMessage](connection types.IConnection, builderFunc func(builder *options.RabbitMQConsumerOptionsBuilder[T]), handler consumer.ConsumerHandler[T], eventSerializer serializer.EventSerializer, logger logger.Logger) (consumer.Consumer, error) {
	builder := options.NewRabbitMQConsumerOptionsBuilder[T]()
	if builderFunc != nil {
		builderFunc(builder)
	}

	// get a new channel on the connection - channel is unique for each consumer
	ch, err := connection.Channel()
	if err != nil {
		return nil, err
	}

	return &RabbitMQConsumer[T]{rabbitmqConsumerOptions: builder.Build(), connection: connection, handler: handler, channel: ch, eventSerializer: eventSerializer, logger: logger}, nil
}

func (r *RabbitMQConsumer[T]) Consume(ctx context.Context) error {
	//https://github.com/rabbitmq/rabbitmq-tutorials/blob/master/go/receive.go
	if r.connection == nil {
		return errors.New("connection is nil")
	}

	if r.connection.IsClosed() {
		return errors.New("connection closed or timeout")
	}

	if err := r.channel.Qos(r.rabbitmqConsumerOptions.PrefetchCount, 0, false); err != nil {
		return err
	}

	err := r.channel.ExchangeDeclare(
		r.rabbitmqConsumerOptions.ExchangeOptions.Name,
		string(r.rabbitmqConsumerOptions.ExchangeOptions.Type),
		r.rabbitmqConsumerOptions.ExchangeOptions.Durable,
		r.rabbitmqConsumerOptions.ExchangeOptions.AutoDelete,
		false,
		r.rabbitmqConsumerOptions.NoWait,
		r.rabbitmqConsumerOptions.ExchangeOptions.Args)
	if err != nil {
		return err
	}

	_, err = r.channel.QueueDeclare(
		r.rabbitmqConsumerOptions.QueueOptions.Name,
		r.rabbitmqConsumerOptions.QueueOptions.Durable,
		r.rabbitmqConsumerOptions.QueueOptions.AutoDelete,
		r.rabbitmqConsumerOptions.QueueOptions.Exclusive,
		r.rabbitmqConsumerOptions.NoWait,
		r.rabbitmqConsumerOptions.QueueOptions.Args)
	if err != nil {
		return err
	}

	err = r.channel.QueueBind(
		r.rabbitmqConsumerOptions.QueueOptions.Name,
		r.rabbitmqConsumerOptions.BindingOptions.RoutingKey,
		r.rabbitmqConsumerOptions.ExchangeOptions.Name,
		r.rabbitmqConsumerOptions.NoWait,
		r.rabbitmqConsumerOptions.BindingOptions.Args)
	if err != nil {
		return err
	}

	deliveryChan, err := r.channel.Consume(
		r.rabbitmqConsumerOptions.QueueOptions.Name,
		r.rabbitmqConsumerOptions.ConsumerId,
		r.rabbitmqConsumerOptions.AutoAck,
		r.rabbitmqConsumerOptions.QueueOptions.Exclusive,
		r.rabbitmqConsumerOptions.NoLocal,
		r.rabbitmqConsumerOptions.NoWait,
		nil)
	if err != nil {
		return err
	}
	go func() {
		for delivery := range deliveryChan {
			//https://github.com/streadway/amqp/blob/2aa28536587a0090d8280eed56c75867ce7e93ec/delivery.go#L62
			r.handleReceived(ctx, delivery, r.handler)
		}
	}()

	return nil
}

func (r *RabbitMQConsumer[T]) UnConsume(ctx context.Context) error {
	err := r.channel.Cancel(r.rabbitmqConsumerOptions.ConsumerId, false)
	if err != nil {
		return err
	}

	defer r.channel.Close() // TODO: this error must be logged

	done := make(chan struct{}, 1)

	go func() {
		for {
			if len(r.deliveryRoutines) == 0 {
				done <- struct{}{}
			}
		}
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (r *RabbitMQConsumer[T]) handleReceived(ctx context.Context, delivery amqp091.Delivery, handler consumer.ConsumerHandler[T]) {
	r.deliveryRoutines <- struct{}{}

	go func() {
		defer func() { <-r.deliveryRoutines }()

		consumeContext := r.createConsumeContext(delivery)

		ack := func() {
			if err := delivery.Ack(false); err != nil {
				// TODO: this error must be logged
				return
			}
		}
		nack := func() {
			if err := delivery.Nack(false, true); err != nil {
				// TODO: this error must be logged
				return
			}
		}
		r.handle(ctx, ack, nack, consumeContext, handler)
	}()
}

func (r *RabbitMQConsumer[T]) handle(ctx context.Context, ack func(), nack func(), messageConsumeContext types2.IMessageConsumeContext[T], handler consumer.ConsumerHandler[T]) {
	err := retry.Do(func() error {
		err := handler.Handle(ctx, messageConsumeContext)
		return err
	}, append(retryOptions, retry.Context(ctx))...)

	if err != nil {
		r.logger.Error("[RabbitMQConsumer.Handle] error in handling consume message of RabbitmqMQ")
		nack()
	}

	ack()
}

func (r *RabbitMQConsumer[T]) createConsumeContext(delivery amqp091.Delivery) types2.IMessageConsumeContext[T] {
	message := r.deserializeData(delivery.ContentType, delivery.Type, delivery.Body)
	var metadata core.Metadata
	if delivery.Headers != nil {
		metadata = core.MapToMetadata(delivery.Headers)
	}
	consumeContext := types2.NewMessageConsumeContext[T](message, metadata, delivery.ContentType, delivery.Type, delivery.Timestamp, delivery.DeliveryTag, delivery.MessageId, delivery.CorrelationId)

	return consumeContext
}

func (r *RabbitMQConsumer[T]) deserializeData(contentType string, eventType string, body []byte) T {
	if contentType == "" {
		contentType = "application/json"
	}
	if body == nil || len(body) == 0 {
		return *new(T)
	}

	deserialize, err := r.eventSerializer.Deserialize(body, eventType, contentType)
	if err != nil {
		return *new(T)
	}

	return deserialize.(T)
}
