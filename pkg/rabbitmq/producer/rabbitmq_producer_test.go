package producer

import (
	"context"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/core/serializer/json"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger/defaultLogger"
	types2 "github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/types"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/producer/options"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/types"
	uuid "github.com/satori/go.uuid"
	"testing"
)

func Test_Publish_Message(t *testing.T) {
	conn, err := types.NewConnection(context.Background(), &config.RabbitMQConfig{
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

	rabbitmqProducer, err := NewRabbitMQProducer(conn, func(builder *options.RabbitMQProducerOptionsBuilder) {
		builder.WithExchangeType(types.ExchangeTopic)
	}, defaultLogger.Logger, json.NewJsonEventSerializer())
	if err != nil {
		t.Fatal(err)
	}

	err = rabbitmqProducer.Publish(context.Background(), "", NewProducerMessage("test"), nil)
	if err != nil {
		return
	}
}

type ProducerMessage struct {
	*types2.Message
	Data string
}

func NewProducerMessage(data string) *ProducerMessage {
	return &ProducerMessage{
		Data:    data,
		Message: types2.NewMessage(uuid.NewV4().String()),
	}
}
