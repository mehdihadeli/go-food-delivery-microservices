package consumer

import (
	"context"
	"emperror.dev/errors"
	"fmt"
	"github.com/avast/retry-go"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/core"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/core/serializer"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/consumer"
	types2 "github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/types"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/consumer/options"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/rabbitmqErrors"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/types"
	"github.com/rabbitmq/amqp091-go"
	"reflect"
	"sync"
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
	deliveryRoutines        chan struct{} // chan should init before using channel
	eventSerializer         serializer.EventSerializer
	logger                  logger.Logger
	wg                      *sync.WaitGroup
}

func NewRabbitMQConsumer[T types2.IMessage](connection types.IConnection, builderFunc func(builder *options.RabbitMQConsumerOptionsBuilder[T]), eventSerializer serializer.EventSerializer, logger logger.Logger, handler consumer.ConsumerHandler[T]) (consumer.Consumer, error) {
	builder := options.NewRabbitMQConsumerOptionsBuilder[T]()
	if builderFunc != nil {
		builderFunc(builder)
	}

	consumerConfig := builder.Build()
	deliveryRoutines := make(chan struct{}, consumerConfig.ConcurrencyLimit)
	wg := &sync.WaitGroup{}

	cons := &RabbitMQConsumer[T]{rabbitmqConsumerOptions: consumerConfig, deliveryRoutines: deliveryRoutines, wg: wg, connection: connection, handler: handler, eventSerializer: eventSerializer, logger: logger}

	return cons, nil
}

func (r *RabbitMQConsumer[T]) Consume(ctx context.Context) error {
	//https://github.com/rabbitmq/rabbitmq-tutorials/blob/master/go/receive.go
	if r.connection == nil {
		return errors.New("connection is nil")
	}

	r.reConsumeOnDropConnection(ctx)

	// get a new channel on the connection - channel is unique for each consumer
	ch, err := r.connection.Channel()
	if err != nil {
		return rabbitmqErrors.ErrDisconnected
	}
	r.channel = ch

	// The prefetch count tells the Rabbit connection how many messages to retrieve from the server per request.
	prefetchCount := r.rabbitmqConsumerOptions.ConcurrencyLimit * r.rabbitmqConsumerOptions.PrefetchCount
	if err := r.channel.Qos(prefetchCount, 0, false); err != nil {
		return err
	}

	err = r.channel.ExchangeDeclare(
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

	msgs, err := r.channel.Consume(
		r.rabbitmqConsumerOptions.QueueOptions.Name,
		r.rabbitmqConsumerOptions.ConsumerId,
		r.rabbitmqConsumerOptions.AutoAck, //When autoAck (also known as noAck) is true, the server will acknowledge deliveries to this consumer prior to writing the delivery to the network. When autoAck is true, the consumer should not call Delivery.Ack.
		r.rabbitmqConsumerOptions.QueueOptions.Exclusive,
		r.rabbitmqConsumerOptions.NoLocal,
		r.rabbitmqConsumerOptions.NoWait,
		nil)
	if err != nil {
		return err
	}

	//https://blog.boot.dev/golang/connecting-to-rabbitmq-in-golang/
	//https://levelup.gitconnected.com/connecting-a-service-in-golang-to-a-rabbitmq-server-835294d8c914
	//https://www.ribice.ba/golang-rabbitmq-client/
	//https://medium.com/@dhanushgopinath/automatically-recovering-rabbitmq-connections-in-go-applications-7795a605ca59
	for i := 0; i < r.rabbitmqConsumerOptions.ConcurrencyLimit; i++ {
		r.logger.Infof("Processing messages on thread %d", i)
		go func() {
			for {
				select {
				case msg, ok := <-msgs:
					if !ok {
						fmt.Println(r.connection.IsClosed())
						fmt.Println(r.connection.IsConnected())

						r.logger.Error("consumer connection dropped")
						return
					}
					//https://github.com/streadway/amqp/blob/2aa28536587a0090d8280eed56c75867ce7e93ec/delivery.go#L62
					r.handleReceived(ctx, msg, r.handler)
				case <-ctx.Done(): // context canceled, it can stop getting new messages
					err := r.UnConsume(ctx)
					if err != nil {
						r.logger.Error("error in canceling consumer")
						return
					}
					r.logger.Error("consumer canceled")
				}
			}
		}()
	}

	return nil
}

func (r *RabbitMQConsumer[T]) UnConsume(ctx context.Context) error {
	if r.channel != nil || r.channel.IsClosed() == false {
		err := r.channel.Cancel(r.rabbitmqConsumerOptions.ConsumerId, false)
		if err != nil {
			return err
		}
	}

	defer func() {
		if r.channel != nil || r.channel.IsClosed() == false {
			r.channel.Close() // TODO: this error must be logged
		}
	}()

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

func (r *RabbitMQConsumer[T]) reConsumeOnDropConnection(ctx context.Context) {
	go func() {
		for {
			select {
			case reconnect := <-r.connection.ReconnectedChannel():
				if reflect.ValueOf(reconnect).IsValid() {
					r.logger.Info("reconsume_on_drop_connection started")
					err := r.Consume(ctx)
					if err != nil {
						r.logger.Error("reconsume_on_drop_connection finished with error: %v", err)
						continue
					}
					r.logger.Info("reconsume_on_drop_connection finished successfully")
					return
				}
			}
		}
	}()
}

func (r *RabbitMQConsumer[T]) handleReceived(ctx context.Context, delivery amqp091.Delivery, handler consumer.ConsumerHandler[T]) {
	// for ensuring our handler execute completely after shutdown
	r.deliveryRoutines <- struct{}{}

	defer func() { <-r.deliveryRoutines }()

	consumeContext := r.createConsumeContext(delivery)

	var ack func()
	var nack func()

	// if auto-ack is enabled we should not call Ack method manually it could create some unexpected errors
	if r.rabbitmqConsumerOptions.AutoAck == false {
		ack = func() {
			if err := delivery.Ack(false); err != nil {
				r.logger.Error("error sending ACK to RabbitMQ consumer: %v", err)
				return
			}
		}

		nack = func() {
			if err := delivery.Nack(false, true); err != nil {
				r.logger.Error("error in sending Nack to RabbitMQ consumer: %v", err)
				return
			}
		}
	}

	r.handle(ctx, ack, nack, consumeContext, handler)
}

func (r *RabbitMQConsumer[T]) handle(ctx context.Context, ack func(), nack func(), messageConsumeContext types2.IMessageConsumeContext[T], handler consumer.ConsumerHandler[T]) {
	err := retry.Do(func() error {
		err := handler.Handle(ctx, messageConsumeContext)
		return err
	}, append(retryOptions, retry.Context(ctx))...)

	if err != nil {
		r.logger.Error("[RabbitMQConsumer.Handle] error in handling consume message of RabbitmqMQ")
		if nack != nil && r.rabbitmqConsumerOptions.AutoAck == false {
			nack()
		}
	} else if err == nil && ack != nil && r.rabbitmqConsumerOptions.AutoAck == false {
		ack()
	}
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

	if contentType == "application/json" {
		deserialize, err := r.eventSerializer.Deserialize(body, eventType, contentType)
		if err != nil {
			return *new(T)
		}

		return deserialize.(T)
	}

	return *new(T)
}
