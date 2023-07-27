package rabbitmq

import (
	"context"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/messaging/bus"
	rabbitmqBus "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/rabbitmq/bus"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/rabbitmq/config"
	rabbitmqConfigurations "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/rabbitmq/configurations"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orders/internal/orders/configurations/rabbitmq"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orders/internal/shared/contracts"
)

func ConfigOrdersRabbitMQ(ctx context.Context, mqConfig *config.RabbitMQConfig, infra *contracts.InfrastructureConfigurations) (bus.Bus, error) {
	return rabbitmqBus.NewRabbitMQBus(
		ctx,
		mqConfig,
		func(builder rabbitmqConfigurations.RabbitMQConfigurationBuilder) {
			// Orders RabbitMQ configuration
			rabbitmq.ConfigOrdersRabbitMQ(builder)
		},
		infra.EventSerializer,
		infra.Log)
}
