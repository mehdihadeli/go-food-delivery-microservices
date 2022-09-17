package consumer

import (
	"context"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger/defaultLogger"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/consumer"
	types2 "github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/types"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/bus"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/consumer/options"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/types"
	"testing"
)

func Test_Consume_Message(t *testing.T) {
	conn, err := types.NewConnection(&config.RabbitMQConfig{
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

	rabbitmqConsumer, err := NewRabbitMQConsumer[testMessage](conn, func(builder *options.RabbitMQConsumerOptionsBuilder[testMessage]) {
	}, NewTestMessageHandler())

	var consumers []consumer.Consumer
	consumers = append(consumers, rabbitmqConsumer)
	b := bus.NewRabbitMQBus(defaultLogger.Logger, consumers)
	err = b.Start(context.Background())
	if err != nil {
		return
	}
}

type testMessage struct {
	*types2.Message
	Data string
}

type TestMessageHandler struct {
}

func (t TestMessageHandler) Handle(ctx context.Context, consumeContext types2.IMessageConsumeContext[testMessage]) error {
	return nil
}

func NewTestMessageHandler() *TestMessageHandler {
	return &TestMessageHandler{}
}
