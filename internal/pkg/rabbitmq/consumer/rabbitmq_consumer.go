//go:build go1.18

package consumer

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"emperror.dev/errors"
	"github.com/ahmetb/go-linq/v3"
	"github.com/avast/retry-go"
	"github.com/rabbitmq/amqp091-go"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/core/metadata"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/core/serializer"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/messaging/consumer"
	consumeTracing "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/messaging/otel/tracing/consumer"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/messaging/pipeline"
	messagingTypes "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/messaging/types"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/messaging/utils"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/rabbitmq/consumer/configurations"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/rabbitmq/rabbitmqErrors"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/rabbitmq/types"
	errorUtils "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/utils/error_utils"
)

const (
	retryAttempts = 3
	retryDelay    = 300 * time.Millisecond
)

var retryOptions = []retry.Option{
	retry.Attempts(retryAttempts),
	retry.Delay(retryDelay),
	retry.DelayType(retry.BackOffDelay),
}

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
	isConsumedNotifications []func(message messagingTypes.IMessage)
}

// NewRabbitMQConsumer create a new generic RabbitMQ consumer
func NewRabbitMQConsumer(
	connection types.IConnection,
	consumerConfiguration *configurations.RabbitMQConsumerConfiguration,
	eventSerializer serializer.EventSerializer,
	logger logger.Logger,
	isConsumedNotifications ...func(message messagingTypes.IMessage),
) (consumer.Consumer, error) {
	if consumerConfiguration == nil {
		return nil, errors.New("consumer configuration is required")
	}

	if consumerConfiguration.ConsumerMessageType == nil {
		return nil, errors.New("consumer ConsumerMessageType property is required")
	}

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

	cons.isConsumedNotifications = isConsumedNotifications

	return cons, nil
}

func (r *rabbitMQConsumer) IsConsumed(h func(message messagingTypes.IMessage)) {
	r.isConsumedNotifications = append(r.isConsumedNotifications, h)
}

func (r *rabbitMQConsumer) Start(ctx context.Context) error {
	// https://github.com/rabbitmq/rabbitmq-tutorials/blob/master/go/receive.go
	if r.connection == nil {
		return errors.New("connection is nil")
	}

	var exchange string
	var queue string
	var routingKey string

	if r.rabbitmqConsumerOptions.ExchangeOptions.Name != "" {
		exchange = r.rabbitmqConsumerOptions.ExchangeOptions.Name
	} else {
		exchange = utils.GetTopicOrExchangeNameFromType(r.rabbitmqConsumerOptions.ConsumerMessageType)
	}

	if r.rabbitmqConsumerOptions.BindingOptions.RoutingKey != "" {
		routingKey = r.rabbitmqConsumerOptions.BindingOptions.RoutingKey
	} else {
		routingKey = utils.GetRoutingKeyFromType(r.rabbitmqConsumerOptions.ConsumerMessageType)
	}

	if r.rabbitmqConsumerOptions.QueueOptions.Name != "" {
		queue = r.rabbitmqConsumerOptions.QueueOptions.Name
	} else {
		queue = utils.GetQueueNameFromType(r.rabbitmqConsumerOptions.ConsumerMessageType)
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
		exchange,
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
		queue,
		r.rabbitmqConsumerOptions.QueueOptions.Durable,
		r.rabbitmqConsumerOptions.QueueOptions.AutoDelete,
		r.rabbitmqConsumerOptions.QueueOptions.Exclusive,
		r.rabbitmqConsumerOptions.NoWait,
		r.rabbitmqConsumerOptions.QueueOptions.Args)
	if err != nil {
		return err
	}

	err = r.channel.QueueBind(
		queue,
		routingKey,
		exchange,
		r.rabbitmqConsumerOptions.NoWait,
		r.rabbitmqConsumerOptions.BindingOptions.Args)
	if err != nil {
		return err
	}

	msgs, err := r.channel.Consume(
		queue,
		r.rabbitmqConsumerOptions.ConsumerId,
		r.rabbitmqConsumerOptions.AutoAck, // When autoAck (also known as noAck) is true, the server will acknowledge deliveries to this consumer prior to writing the delivery to the network. When autoAck is true, the consumer should not call Delivery.Ack.
		r.rabbitmqConsumerOptions.QueueOptions.Exclusive,
		r.rabbitmqConsumerOptions.NoLocal,
		r.rabbitmqConsumerOptions.NoWait,
		nil,
	)
	if err != nil {
		return err
	}

	// This channel will receive a notification when a channel closed event happens.
	// https://github.com/streadway/amqp/blob/v1.0.0/channel.go#L447
	// https://github.com/rabbitmq/amqp091-go/blob/main/example_client_test.go#L75
	chClosedCh := make(chan *amqp091.Error, 1)
	ch.NotifyClose(chClosedCh)

	// https://blog.boot.dev/golang/connecting-to-rabbitmq-in-golang/
	// https://levelup.gitconnected.com/connecting-a-service-in-golang-to-a-rabbitmq-server-835294d8c914
	// https://www.ribice.ba/golang-rabbitmq-client/
	// https://medium.com/@dhanushgopinath/automatically-recovering-rabbitmq-connections-in-go-applications-7795a605ca59
	// https://github.com/rabbitmq/amqp091-go/blob/main/_examples/pubsub/pubsub.go
	for i := 0; i < r.rabbitmqConsumerOptions.ConcurrencyLimit; i++ {
		r.logger.Infof("Processing messages on thread %d", i)
		go func() {
			for {
				select {
				case <-ctx.Done():
					r.logger.Info("shutting down consumer")
					return
				case amqErr := <-chClosedCh:
					// This case handles the event of closed channel e.g. abnormal shutdown
					r.logger.Errorf("AMQP Channel closed due to: %s", amqErr)

					// Re-set channel to receive notifications
					chClosedCh = make(chan *amqp091.Error, 1)
					ch.NotifyClose(chClosedCh)
				case msg, ok := <-msgs:
					if !ok {
						r.logger.Error("consumer connection dropped")
						return
					}

					// handle received message and remove message form queue with a manual ack
					r.handleReceived(ctx, msg)
				}
			}
		}()
	}

	return nil
}

func (r *rabbitMQConsumer) Stop() error {
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

func (r *rabbitMQConsumer) ConnectHandler(handler consumer.ConsumerHandler) {
	r.handlers = append(r.handlers, handler)
}

func (r *rabbitMQConsumer) GetName() string {
	return r.rabbitmqConsumerOptions.Name
}

func (r *rabbitMQConsumer) reConsumeOnDropConnection(ctx context.Context) {
	go func() {
		defer errorUtils.HandlePanic()
		for {
			select {
			case reconnect := <-r.connection.ReconnectedChannel():
				if reflect.ValueOf(reconnect).IsValid() {
					r.logger.Info("reconsume_on_drop_connection started")
					err := r.Start(ctx)
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
			semconv.MessagingRabbitmqDestinationRoutingKey(delivery.RoutingKey),
		},
	}
	ctx, beforeConsumeSpan := consumeTracing.StartConsumerSpan(
		ctx,
		&meta,
		string(delivery.Body),
		consumerTraceOption,
	)

	consumeContext, err := r.createConsumeContext(delivery)
	if err != nil {
		r.logger.Error(consumeTracing.FinishConsumerSpan(beforeConsumeSpan, err))
		return
	}

	var ack func()
	var nack func()

	// if auto-ack is enabled we should not call Ack method manually it could create some unexpected errors
	if r.rabbitmqConsumerOptions.AutoAck == false {
		ack = func() {
			if err := delivery.Ack(false); err != nil {
				r.logger.Error(
					"error sending ACK to RabbitMQ consumer: %v",
					consumeTracing.FinishConsumerSpan(beforeConsumeSpan, err),
				)
				return
			}
			_ = consumeTracing.FinishConsumerSpan(beforeConsumeSpan, nil)
			if len(r.isConsumedNotifications) > 0 {
				for _, notification := range r.isConsumedNotifications {
					if notification != nil {
						notification(consumeContext.Message())
					}
				}
			}
		}

		nack = func() {
			if err := delivery.Nack(false, true); err != nil {
				r.logger.Error(
					"error in sending Nack to RabbitMQ consumer: %v",
					consumeTracing.FinishConsumerSpan(beforeConsumeSpan, err),
				)
				return
			}
			_ = consumeTracing.FinishConsumerSpan(beforeConsumeSpan, nil)
		}
	}

	r.handle(ctx, ack, nack, consumeContext)
}

func (r *rabbitMQConsumer) handle(
	ctx context.Context,
	ack func(),
	nack func(),
	messageConsumeContext messagingTypes.MessageConsumeContext,
) {
	var err error
	for _, handler := range r.handlers {
		err = r.runHandlersWithRetry(ctx, handler, messageConsumeContext)
		if err != nil {
			break
		}
	}

	if err != nil {
		r.logger.Error(
			"[rabbitMQConsumer.Handle] error in handling consume message of RabbitmqMQ, prepare for nacking message",
		)
		if nack != nil && r.rabbitmqConsumerOptions.AutoAck == false {
			nack()
		}
	} else if err == nil && ack != nil && r.rabbitmqConsumerOptions.AutoAck == false {
		ack()
	}
}

func (r *rabbitMQConsumer) runHandlersWithRetry(
	ctx context.Context,
	handler consumer.ConsumerHandler,
	messageConsumeContext messagingTypes.MessageConsumeContext,
) error {
	err := retry.Do(func() error {
		var lastHandler pipeline.ConsumerHandlerFunc

		if r.pipelines != nil && len(r.pipelines) > 0 {
			reversPipes := r.reversOrder(r.pipelines)
			lastHandler = func() error {
				handler := handler.(consumer.ConsumerHandler)
				return handler.Handle(ctx, messageConsumeContext)
			}

			aggregateResult := linq.From(reversPipes).
				AggregateWithSeedT(lastHandler, func(next pipeline.ConsumerHandlerFunc, pipe pipeline.ConsumerPipeline) pipeline.ConsumerHandlerFunc {
					pipeValue := pipe
					nexValue := next

					var handlerFunc pipeline.ConsumerHandlerFunc = func() error {
						genericContext, ok := messageConsumeContext.(messagingTypes.MessageConsumeContext)
						if ok {
							return pipeValue.Handle(ctx, genericContext, nexValue)
						}
						return pipeValue.Handle(
							ctx,
							messageConsumeContext.(messagingTypes.MessageConsumeContext),
							nexValue,
						)
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

	return err
}

func (r *rabbitMQConsumer) createConsumeContext(
	delivery amqp091.Delivery,
) (messagingTypes.MessageConsumeContext, error) {
	message := r.deserializeData(delivery.ContentType, delivery.Type, delivery.Body)
	if reflect.ValueOf(message).IsZero() || reflect.ValueOf(message).IsNil() {
		return *new(messagingTypes.MessageConsumeContext), errors.New(
			"error in deserialization of payload",
		)
	}
	m, ok := message.(messagingTypes.IMessage)
	if !ok || m.IsMessage() == false {
		return nil, errors.New(
			fmt.Sprintf(
				"message %s is not a message type or message property is nil",
				utils.GetMessageBaseReflectType(message),
			),
		)
	}

	var meta metadata.Metadata
	if delivery.Headers != nil {
		meta = metadata.MapToMetadata(delivery.Headers)
	}

	consumeContext := messagingTypes.NewMessageConsumeContext(
		message.(messagingTypes.IMessage),
		meta,
		delivery.ContentType,
		delivery.Type,
		delivery.Timestamp,
		delivery.DeliveryTag,
		delivery.MessageId,
		delivery.CorrelationId,
	)
	return consumeContext, nil
}

func (r *rabbitMQConsumer) deserializeData(
	contentType string,
	eventType string,
	body []byte,
) interface{} {
	if contentType == "" {
		contentType = "application/json"
	}

	if body == nil || len(body) == 0 {
		r.logger.Error("message body is nil or empty in the consumer")
		return nil
	}

	if contentType == "application/json" {
		// r.rabbitmqConsumerOptions.ConsumerMessageType --> actual type
		// deserialize, err := r.eventSerializer.DeserializeType(body, r.rabbitmqConsumerOptions.ConsumerMessageType, contentType)
		deserialize, err := r.eventSerializer.DeserializeMessage(
			body,
			eventType,
			contentType,
		) // or this to explicit type deserialization
		if err != nil {
			r.logger.Errorf(
				fmt.Sprintf("error in deserilizng of type '%s' in the consumer", eventType),
			)
			return nil
		}

		return deserialize
	}

	return nil
}

func (r *rabbitMQConsumer) reversOrder(
	values []pipeline.ConsumerPipeline,
) []pipeline.ConsumerPipeline {
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
