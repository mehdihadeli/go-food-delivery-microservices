package options

import "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/rabbitmq/types"

type RabbitMQExchangeOptions struct {
	Name       string
	Type       types.ExchangeType
	AutoDelete bool
	Durable    bool
	Args       map[string]any
}
