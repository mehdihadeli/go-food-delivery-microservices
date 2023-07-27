package rabbitmq

import (
	"context"
	"testing"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/rabbitmq/config"
)

var RabbitmqContainerOptionsDecorator = func(t *testing.T, ctx context.Context) interface{} {
	return func(c *config.RabbitmqOptions) (*config.RabbitmqOptions, error) {
		return NewRabbitMQTestContainers().CreatingContainerOptions(ctx, t)
	}
}
