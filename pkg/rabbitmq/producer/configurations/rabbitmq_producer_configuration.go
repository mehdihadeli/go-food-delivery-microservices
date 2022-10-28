//go:build go1.18

package configurations

import (
	types2 "github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/types"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/utils"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/producer/options"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/types"
	"reflect"
)

type RabbitMQProducerConfiguration struct {
	ProducerMessageType reflect.Type
	ExchangeOptions     *options.RabbitMQExchangeOptions
	RoutingKey          string
	DeliveryMode        uint8
	Priority            uint8
	AppId               string
	Expiration          string
	ReplyTo             string
	ContentEncoding     string
}

func NewDefaultRabbitMQProducerConfiguration(messageType types2.IMessage) *RabbitMQProducerConfiguration {
	return &RabbitMQProducerConfiguration{
		ExchangeOptions:     &options.RabbitMQExchangeOptions{Durable: true, Type: types.ExchangeTopic, Name: utils.GetTopicOrExchangeName(messageType)},
		DeliveryMode:        2,
		RoutingKey:          utils.GetRoutingKey(messageType),
		ProducerMessageType: utils.GetMessageBaseReflectType(messageType),
	}
}
