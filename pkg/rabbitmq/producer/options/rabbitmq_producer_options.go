package options

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/types"
)

type RabbitMQProducerOptions struct {
	ExchangeOptions *RabbitMQExchangeOptions
}

func NewDefaultRabbitMQProducerOptions() *RabbitMQProducerOptions {
	return &RabbitMQProducerOptions{
		ExchangeOptions: &RabbitMQExchangeOptions{Durable: true, Type: types.ExchangeTopic},
	}
}
