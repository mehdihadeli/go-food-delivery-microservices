package consumer

import (
	"context"
	"fmt"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/core/serializer/json"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger/defaultLogger"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/consumer"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/pipeline"
	types2 "github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/types"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/bus"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/consumer/options"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/producer"
	options2 "github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/producer/options"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/types"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/test"
	uuid "github.com/satori/go.uuid"
	"testing"
	"time"
)

func Test_Consume_Message(t *testing.T) {
	test.SkipCI(t)
	conn, err := types.NewRabbitMQConnection(context.Background(), &config.RabbitMQConfig{
		RabbitMqHostOptions: &config.RabbitMqHostOptions{
			UserName: "guest",
			Password: "guest",
			HostName: "localhost",
			Port:     5672,
		},
	})
	if err != nil {
		t.Fatal(err)
		return
	}

	rabbitmqConsumer, err := NewRabbitMQConsumer[*ProducerConsumerMessage](
		json.NewJsonEventSerializer(),
		defaultLogger.Logger,
		conn,
		func(builder *options.RabbitMQConsumerOptionsBuilder[*ProducerConsumerMessage]) {
			//builder.WithAutoAck(true)
		},
		NewTestMessageHandler(),
		NewPipeline1[*ProducerConsumerMessage]())

	var consumers []consumer.Consumer
	consumers = append(consumers, rabbitmqConsumer)

	b := bus.NewRabbitMQBus(defaultLogger.Logger, consumers)
	err = b.Start(context.Background())
	if err != nil {
		return
	}

	rabbitmqProducer, err := producer.NewRabbitMQProducer(
		conn,
		func(builder *options2.RabbitMQProducerOptionsBuilder) {
			builder.WithExchangeType(types.ExchangeTopic)
		},
		defaultLogger.Logger,
		json.NewJsonEventSerializer())
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second * 5)

	fmt.Println("closing connection")
	conn.Close()
	fmt.Println(conn.IsClosed())

	time.Sleep(time.Second * 10)
	fmt.Println("after 10 second of closing connection")
	fmt.Println(conn.IsClosed())

	err = rabbitmqProducer.Publish(context.Background(), NewProducerConsumerMessage("test"), nil)
	for err != nil {
		err = rabbitmqProducer.Publish(context.Background(), NewProducerConsumerMessage("test"), nil)
	}

	err = rabbitmqProducer.Publish(context.Background(), NewProducerConsumerMessage("test"), nil)
	for err != nil {
		err = rabbitmqProducer.Publish(context.Background(), NewProducerConsumerMessage("test"), nil)
	}

	time.Sleep(time.Second * 5)
	fmt.Println(conn.IsClosed())
	fmt.Println(conn.IsConnected())
}

type ProducerConsumerMessage struct {
	*types2.Message
	Data string
}

func NewProducerConsumerMessage(data string) *ProducerConsumerMessage {
	return &ProducerConsumerMessage{
		Data:    data,
		Message: types2.NewMessage(uuid.NewV4().String()),
	}
}

// /////////// ConsumerHandler
type TestMessageHandler struct {
}

func (t *TestMessageHandler) Handle(ctx context.Context, consumeContext types2.IMessageConsumeContext[*ProducerConsumerMessage]) error {
	message := consumeContext.Message()
	fmt.Println(message)

	return nil
}

func NewTestMessageHandler() *TestMessageHandler {
	return &TestMessageHandler{}
}

// /////////////// ConsumerPipeline
type Pipeline1[T types2.IMessage] struct {
}

func NewPipeline1[T types2.IMessage]() pipeline.ConsumerPipeline[T] {
	return &Pipeline1[T]{}
}
func (p Pipeline1[T]) Handle(ctx context.Context, consumerContext types2.IMessageConsumeContext[T], next pipeline.ConsumerHandlerFunc) error {
	fmt.Println("PipelineBehaviourTest.Handled")

	fmt.Println(fmt.Sprintf("pipeline got a message with id '%s'", consumerContext.Message().GeMessageId()))

	err := next()
	if err != nil {
		return err
	}
	return nil
}
