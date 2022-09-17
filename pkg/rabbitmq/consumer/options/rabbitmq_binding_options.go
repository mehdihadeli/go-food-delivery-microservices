package options

type RabbitMQBindingOptions struct {
	RoutingKey string
	Args       map[string]any
}
