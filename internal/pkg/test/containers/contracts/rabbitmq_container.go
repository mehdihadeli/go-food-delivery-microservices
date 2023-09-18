package contracts

import (
	"context"
	"testing"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/core/serializer"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/messaging/bus"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/rabbitmq/config"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/rabbitmq/configurations"
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

type RabbitMQContainer interface {
	Start(ctx context.Context,
		t *testing.T,
		serializer serializer.EventSerializer,
		rabbitmqBuilderFunc configurations.RabbitMQConfigurationBuilderFuc,
		options ...*RabbitMQContainerOptions) (bus.Bus, error)

	CreatingContainerOptions(
		ctx context.Context,
		t *testing.T,
		options ...*RabbitMQContainerOptions,
	) (*config.RabbitmqHostOptions, error)

	Cleanup(ctx context.Context) error
}
