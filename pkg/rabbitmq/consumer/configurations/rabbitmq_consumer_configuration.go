package configurations

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/consumer"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/pipeline"
	types2 "github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/types"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/utils"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/consumer/options"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/types"
	"reflect"
)

type RabbitMQConsumerConfiguration struct {
	ConsumerMessageType reflect.Type
	Pipelines           []pipeline.ConsumerPipeline
	Handlers            []consumer.ConsumerHandler
	*consumer.ConsumerOptions
	ConcurrencyLimit int
	// The prefetch count tells the Rabbit connection how many messages to retrieve from the server per request.
	PrefetchCount   int
	AutoAck         bool
	NoLocal         bool
	NoWait          bool
	BindingOptions  *options.RabbitMQBindingOptions
	QueueOptions    *options.RabbitMQQueueOptions
	ExchangeOptions *options.RabbitMQExchangeOptions
}

func NewDefaultRabbitMQConsumerConfiguration() *RabbitMQConsumerConfiguration {
	return &RabbitMQConsumerConfiguration{
		ConsumerOptions:  &consumer.ConsumerOptions{ExitOnError: false, ConsumerId: ""},
		ConcurrencyLimit: 1,
		PrefetchCount:    4, //how many messages we can handle at once
		NoLocal:          false,
		NoWait:           true,
		BindingOptions:   &options.RabbitMQBindingOptions{},
		ExchangeOptions:  &options.RabbitMQExchangeOptions{Durable: true, Type: types.ExchangeTopic},
		QueueOptions:     &options.RabbitMQQueueOptions{Durable: true},
	}
}

func NewDefaultRabbitMQConsumerConfigurationT(messageType types2.IMessage) *RabbitMQConsumerConfiguration {
	return &RabbitMQConsumerConfiguration{
		ConsumerOptions:  &consumer.ConsumerOptions{ExitOnError: false, ConsumerId: ""},
		ConcurrencyLimit: 1,
		PrefetchCount:    4, //how many messages we can handle at once
		NoLocal:          false,
		NoWait:           true,
		BindingOptions:   &options.RabbitMQBindingOptions{RoutingKey: utils.GetRoutingKey(messageType)},
		ExchangeOptions:  &options.RabbitMQExchangeOptions{Durable: true, Type: types.ExchangeTopic, Name: utils.GetTopicOrExchangeName(messageType)},
		QueueOptions:     &options.RabbitMQQueueOptions{Durable: true, Name: utils.GetQueueName(messageType)},
	}
}
