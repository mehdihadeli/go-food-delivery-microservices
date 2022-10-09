package rabbitmq

import (
	"context"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/bus"
	rabbitmqBus "github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/bus"
	config2 "github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/config"
	rabbitmqConfigurations "github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/configurations"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/configurations/rabbitmq"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/shared/contracts"
)

func ConfigOrdersRabbitMQ(ctx context.Context, mqConfig *config2.RabbitMQConfig, infra contracts.InfrastructureConfigurations) (bus.Bus, error) {
	return rabbitmqBus.NewRabbitMQBus(
		ctx,
		mqConfig,
		func(builder rabbitmqConfigurations.RabbitMQConfigurationBuilder) {
			// Orders RabbitMQ configuration
			rabbitmq.ConfigOrdersRabbitMQ(builder)
		},
		infra.EventSerializer(),
		infra.Log())
}
