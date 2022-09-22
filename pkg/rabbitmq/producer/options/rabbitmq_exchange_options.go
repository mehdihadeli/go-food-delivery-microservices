package options

import "github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/types"

type RabbitMQExchangeOptions struct {
	Type       types.ExchangeType
	AutoDelete bool
	Durable    bool
	Args       map[string]any
}
