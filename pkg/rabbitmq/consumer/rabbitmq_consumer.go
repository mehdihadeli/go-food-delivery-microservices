package consumer

import (
	"context"
	"emperror.dev/errors"
	"fmt"
	"github.com/ahmetb/go-linq/v3"
	"github.com/avast/retry-go"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/core/metadata"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/core/serializer"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/consumer"
	consumeTracing "github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/otel/tracing/consumer"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/pipeline"
	messagingTypes "github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/types"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/consumer/configurations"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/rabbitmqErrors"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/types"
	"github.com/rabbitmq/amqp091-go"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"reflect"
	"time"
)

const (
	retryAttempts = 3
	retryDelay    = 300 * time.Millisecond
)

var (
	retryOptions = []retry.Option{retry.Attempts(retryAttempts), retry.Delay(retryDelay), retry.DelayType(retry.BackOffDelay)}
)

type rabbitMQConsumer struct {
	rabbitmqConsumerOptions *configurations.RabbitMQConsumerConfiguration
	connection              types.IConnection
	handlerDefault          consumer.ConsumerHandler
	channel                 *amqp091.Channel
	deliveryRoutines        chan struct{} // chan should init before using channel
	eventSerializer         serializer.EventSerializer
	logger                  logger.Logger
	ErrChan                 chan error
	handlers                []consumer.ConsumerHandler
	pipelines               []pipeline.ConsumerPipeline
}

// NewRabbitMQConsumer create a new generic RabbitMQ consumer
func NewRabbitMQConsumer(consumerMessageType reflect.Type, eventSerializer serializer.EventSerializer, logger logger.Logger, connection types.IConnection, consumerConfiguration *configurations.RabbitMQConsumerConfiguration) (consumer.Consumer, error) {
	deliveryRoutines := make(chan struct{}, consumerConfiguration.ConcurrencyLimit)

	cons := &rabbitMQConsumer{
		eventSerializer:         eventSerializer,
		logger:                  logger,
		rabbitmqConsumerOptions: consumerConfiguration,
		deliveryRoutines:        deliveryRoutines,
		ErrChan:                 make(chan error),
		connection:              connection,
		handlers:                consumerConfiguration.Handlers,
		pipelines:               consumerConfiguration.Pipelines,
	}

	return cons, nil
}

func (r *rabbitMQConsumer) Consume(ctx context.Context) error {
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
						r.logger.Error("consumer connection dropped")
						return
					}

					//https://github.com/streadway/amqp/blob/2aa28536587a0090d8280eed56c75867ce7e93ec/delivery.go#L62
					r.handleReceived(ctx, msg)
				}
			}
		}()
	}

	return nil
}

func (r *rabbitMQConsumer) UnConsume(ctx context.Context) error {
	defer func() {
		if r.channel != nil && r.channel.IsClosed() == false {
			r.channel.Cancel(r.rabbitmqConsumerOptions.ConsumerId, false)
			r.channel.Close()
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
	}
}

func (r *rabbitMQConsumer) reConsumeOnDropConnection(ctx context.Context) {
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

func (r *rabbitMQConsumer) handleReceived(ctx context.Context, delivery amqp091.Delivery) {
	// for ensuring our handlers execute completely after shutdown
	r.deliveryRoutines <- struct{}{}

	defer func() { <-r.deliveryRoutines }()

	var meta metadata.Metadata
	if delivery.Headers != nil {
		meta = metadata.MapToMetadata(delivery.Headers)
	}

	consumerTraceOption := &consumeTracing.ConsumerTracingOptions{
		MessagingSystem: "rabbitmq",
		DestinationKind: "queue",
		Destination:     r.rabbitmqConsumerOptions.QueueOptions.Name,
		OtherAttributes: []attribute.KeyValue{
			semconv.MessagingRabbitmqRoutingKeyKey.String(delivery.RoutingKey),
		},
	}
	ctx, beforeConsumeSpan := consumeTracing.StartConsumerSpan(ctx, &meta, string(delivery.Body), consumerTraceOption)

	consumeContext := r.createConsumeContext(delivery)
	if consumeContext == nil {
		r.logger.Error(consumeTracing.FinishConsumerSpan(beforeConsumeSpan, errors.New("createConsumeContext is nil")).Error())
		return
	}

	var ack func()
	var nack func()

	// if auto-ack is enabled we should not call Ack method manually it could create some unexpected errors
	if r.rabbitmqConsumerOptions.AutoAck == false {
		ack = func() {
			if err := delivery.Ack(false); err != nil {
				r.logger.Error("error sending ACK to RabbitMQ consumer: %v", consumeTracing.FinishConsumerSpan(beforeConsumeSpan, err))
				return
			}
			_ = consumeTracing.FinishConsumerSpan(beforeConsumeSpan, nil)
		}

		nack = func() {
			if err := delivery.Nack(false, true); err != nil {
				r.logger.Error("error in sending Nack to RabbitMQ consumer: %v", consumeTracing.FinishConsumerSpan(beforeConsumeSpan, err))
				return
			}
			_ = consumeTracing.FinishConsumerSpan(beforeConsumeSpan, nil)
		}
	}

	r.handle(ctx, ack, nack, consumeContext)
}

func (r *rabbitMQConsumer) handle(ctx context.Context, ack func(), nack func(), messageConsumeContext messagingTypes.MessageConsumeContextBase) {
	for _, handler := range r.handlers {
		err := retry.Do(func() error {
			var lastHandler pipeline.ConsumerHandlerFunc

			if r.pipelines != nil && len(r.pipelines) > 0 {
				var reversPipes = r.reversOrder(r.pipelines)
				lastHandler = func() error {
					handler := handler.(consumer.ConsumerHandler)
					return handler.Handle(ctx, messageConsumeContext.(messagingTypes.MessageConsumeContext))
				}

				aggregateResult := linq.From(reversPipes).AggregateWithSeedT(lastHandler, func(next pipeline.ConsumerHandlerFunc, pipe pipeline.ConsumerPipeline) pipeline.ConsumerHandlerFunc {
					pipeValue := pipe
					nexValue := next

					var handlerFunc pipeline.ConsumerHandlerFunc = func() error {
						genericContext, ok := messageConsumeContext.(messagingTypes.MessageConsumeContext)
						if ok {
							return pipeValue.Handle(ctx, genericContext, nexValue)
						}
						return pipeValue.Handle(ctx, messageConsumeContext.(messagingTypes.MessageConsumeContext), nexValue)
					}
					return handlerFunc
				})

				v := aggregateResult.(pipeline.ConsumerHandlerFunc)
				err := v()

				if err != nil {
					return errors.Wrap(err, "error handling consumer handlers pipeline")
				}
				return nil
			} else {
				err := handler.Handle(ctx, messageConsumeContext.(messagingTypes.MessageConsumeContext))
				if err != nil {
					return err
				}
			}
			return nil
		}, append(retryOptions, retry.Context(ctx))...)
		if err != nil {
			r.logger.Error("[rabbitMQConsumer.Handle] error in handling consume message of RabbitmqMQ, prepare for nacking message")
			if nack != nil && r.rabbitmqConsumerOptions.AutoAck == false {
				nack()
			}
		} else if err == nil && ack != nil && r.rabbitmqConsumerOptions.AutoAck == false {
			ack()
		}
	}
}

func (r *rabbitMQConsumer) createConsumeContext(delivery amqp091.Delivery) messagingTypes.MessageConsumeContextBase {
	message := r.deserializeData(delivery.ContentType, delivery.Type, delivery.Body)
	if reflect.ValueOf(message).IsZero() || reflect.ValueOf(message).IsNil() {
		r.logger.Error("error in deserialization of payload")
		return *new(messagingTypes.MessageConsumeContextBase)
	}

	var meta metadata.Metadata
	if delivery.Headers != nil {
		meta = metadata.MapToMetadata(delivery.Headers)
	}

	consumeContext := messagingTypes.NewMessageConsumeContext(message.(messagingTypes.IMessage), meta, delivery.ContentType, delivery.Type, delivery.Timestamp, delivery.DeliveryTag, delivery.MessageId, delivery.CorrelationId)
	return consumeContext
}

func (r *rabbitMQConsumer) deserializeData(contentType string, eventType string, body []byte) interface{} {
	if contentType == "" {
		contentType = "application/json"
	}

	if body == nil || len(body) == 0 {
		r.logger.Error("message body is nil or empty in the consumer")
		return nil
	}

	if contentType == "application/json" {
		// r.rabbitmqConsumerOptions.ConsumerMessageType --> actual type
		//deserialize, err := r.eventSerializer.DeserializeType(body, r.rabbitmqConsumerOptions.ConsumerMessageType, contentType)
		deserialize, err := r.eventSerializer.DeserializeMessage(body, eventType, contentType) // or this to explicit type deserialization
		if err != nil {
			r.logger.Errorf(fmt.Sprintf("error in deserilizng of type '%s' in the consumer", eventType))
			return nil
		}

		return deserialize
	}

	return nil
}

func (r *rabbitMQConsumer) reversOrder(values []pipeline.ConsumerPipeline) []pipeline.ConsumerPipeline {
	var reverseValues []pipeline.ConsumerPipeline

	for i := len(values) - 1; i >= 0; i-- {
		reverseValues = append(reverseValues, values[i])
	}

	return reverseValues
}

func (r *rabbitMQConsumer) existsPipeType(p reflect.Type) bool {
	for _, pipe := range r.pipelines {
		if reflect.TypeOf(pipe) == p {
			return true
		}
	}

	return false
}
