package options

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/types"
)

type RabbitMQProducerOptions struct {
	ExchangeOptions *RabbitMQExchangeOptions
	DeliveryMode    uint8
	Priority        uint8
	AppId           string
	Expiration      string
	ReplyTo         string
	ContentEncoding string
}

func NewDefaultRabbitMQProducerOptions() *RabbitMQProducerOptions {
	return &RabbitMQProducerOptions{
		ExchangeOptions: &RabbitMQExchangeOptions{Durable: true, Type: types.ExchangeTopic},
		DeliveryMode:    2,
	}
}
