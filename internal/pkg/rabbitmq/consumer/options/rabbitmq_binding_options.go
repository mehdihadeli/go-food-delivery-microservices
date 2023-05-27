//go:build go1.18

package options

type RabbitMQBindingOptions struct {
	RoutingKey string
	Args       map[string]any
}
