package rabbitmq

import (
	"context"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/bus"
	rabbitmqBus "github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/bus"
	config2 "github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/config"
	rabbitmqConfigurations "github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/configurations"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/configurations/rabbitmq"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/shared/contracts"
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
