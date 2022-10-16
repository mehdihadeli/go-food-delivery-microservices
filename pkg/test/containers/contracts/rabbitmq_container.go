package contracts

import (
	"context"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/bus"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/configurations"
	"testing"
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
