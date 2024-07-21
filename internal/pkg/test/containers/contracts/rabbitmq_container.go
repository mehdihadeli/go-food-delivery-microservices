package contracts

import (
	"context"
	"fmt"
	"testing"

	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/rabbitmq/config"
)

type RabbitMQContainerOptions struct {
	Host        string
	VirtualHost string
	Ports       []string
	HostPort    int
	HttpPort    int
	UserName    string
	Password    string
	ImageName   string
	Name        string
	Tag         string
}

func (h *RabbitMQContainerOptions) AmqpEndPoint() string {
	return fmt.Sprintf(
		"amqp://%s:%s@%s:%d",
		h.UserName,
		h.Password,
		h.Host,
		h.HostPort,
	)
}

func (h *RabbitMQContainerOptions) HttpEndPoint() string {
	return fmt.Sprintf("http://%s:%d", h.Host, h.HttpPort)
}

type RabbitMQContainer interface {
	PopulateContainerOptions(
		ctx context.Context,
		t *testing.T,
		options ...*RabbitMQContainerOptions,
	) (*config.RabbitmqHostOptions, error)

	Cleanup(ctx context.Context) error
}
