package configurations

import (
	"fmt"
	"reflect"

	consumer2 "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/core/messaging/consumer"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/core/messaging/pipeline"
	types2 "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/core/messaging/types"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/core/messaging/utils"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/rabbitmq/consumer/options"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/rabbitmq/types"
)

type RabbitMQConsumerConfiguration struct {
	Name                string
	ConsumerMessageType reflect.Type
	Pipelines           []pipeline.ConsumerPipeline
	Handlers            []consumer2.ConsumerHandler
	*consumer2.ConsumerOptions
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

func NewDefaultRabbitMQConsumerConfiguration(
	messageType types2.IMessage,
) *RabbitMQConsumerConfiguration {
	name := fmt.Sprintf("%s_consumer", utils.GetMessageName(messageType))

	return &RabbitMQConsumerConfiguration{
		ConsumerOptions:  &consumer2.ConsumerOptions{ExitOnError: false, ConsumerId: ""},
		ConcurrencyLimit: 1,
		PrefetchCount:    4, // how many messages we can handle at once
		NoLocal:          false,
		NoWait:           true,
		BindingOptions: &options.RabbitMQBindingOptions{
			RoutingKey: utils.GetRoutingKey(messageType),
		},
		ExchangeOptions: &options.RabbitMQExchangeOptions{
			Durable: true,
			Type:    types.ExchangeTopic,
			Name:    utils.GetTopicOrExchangeName(messageType),
		},
		QueueOptions: &options.RabbitMQQueueOptions{
			Durable: true,
			Name:    utils.GetQueueName(messageType),
		},
		ConsumerMessageType: utils.GetMessageBaseReflectType(messageType),
		Name:                name,
	}
}
