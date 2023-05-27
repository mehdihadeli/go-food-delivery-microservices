package contracts

import (
    "context"
    "testing"

    "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/messaging/bus"
    "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/rabbitmq/configurations"
)

type RabbitMQContainerOptions struct {
	Host        string
	VirtualHost string
	Ports       []string
	HostPort    int
	UserName    string
	Password    string
	ImageName   string
	Name        string
	Tag         string
}

type RabbitMQContainer interface {
	Start(ctx context.Context, t *testing.T, rabbitmqBuilderFunc configurations.RabbitMQConfigurationBuilderFuc, options ...*RabbitMQContainerOptions) (bus.Bus, error)
	Cleanup(ctx context.Context) error
}
