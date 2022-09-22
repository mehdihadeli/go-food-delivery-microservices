package options

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/consumer"
	types2 "github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/types"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/utils"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/types"
)

type RabbitMQConsumerOptions struct {
	*consumer.ConsumerOptions
	ConcurrencyLimit int
	// The prefetch count tells the Rabbit connection how many messages to retrieve from the server per request.
	PrefetchCount   int
	AutoAck         bool
	NoLocal         bool
	NoWait          bool
	BindingOptions  *RabbitMQBindingOptions
	QueueOptions    *RabbitMQQueueOptions
	ExchangeOptions *RabbitMQExchangeOptions
}

func NewDefaultRabbitMQConsumerOptions[T types2.IMessage]() *RabbitMQConsumerOptions {
	return &RabbitMQConsumerOptions{
		ConsumerOptions:  &consumer.ConsumerOptions{ExitOnError: false, ConsumerId: ""},
		ConcurrencyLimit: 1,
		PrefetchCount:    4, //how many messages we can handle at once
		NoLocal:          false,
		NoWait:           true,
		BindingOptions:   &RabbitMQBindingOptions{RoutingKey: utils.GetRoutingKey(*new(T))},
		ExchangeOptions:  &RabbitMQExchangeOptions{Durable: true, Type: types.ExchangeTopic, Name: utils.GetTopicOrExchangeName(*new(T))},
		QueueOptions:     &RabbitMQQueueOptions{Durable: true, Name: utils.GetQueueName(*new(T))},
	}
}
