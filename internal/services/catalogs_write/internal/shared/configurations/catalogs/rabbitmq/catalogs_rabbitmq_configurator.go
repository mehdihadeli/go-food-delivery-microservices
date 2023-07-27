package rabbitmq

import (
	"context"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/messaging/bus"
	rabbitmqBus "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/rabbitmq/bus"
	config2 "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/rabbitmq/config"
	rabbitmqConfigurations "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/rabbitmq/configurations"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/write_service/internal/products/configurations/rabbitmq"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/write_service/internal/shared/contracts"
)

func ConfigCatalogsRabbitMQ(ctx context.Context, mqConfig *config2.RabbitMQConfig, infra *contracts.InfrastructureConfigurations) (bus.Bus, error) {
	return rabbitmqBus.NewRabbitMQBus(
		ctx,
		mqConfig,
		func(builder rabbitmqConfigurations.RabbitMQConfigurationBuilder) {
			// Products RabbitMQ configuration
			rabbitmq.ConfigProductsRabbitMQ(builder)
		},
		infra.EventSerializer,
		infra.Log)
}
