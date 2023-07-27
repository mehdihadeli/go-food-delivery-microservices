//go:build.sh go1.18

package options

type RabbitMQQueueOptions struct {
	Name       string
	Durable    bool
	Exclusive  bool
	AutoDelete bool
	Args       map[string]any
}
