package rabbitmq

import (
	rabbitmqConfigurations "github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/configurations"
	producerConfigurations "github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/producer/configurations"
	v1 "github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/features/creating_order/events/integration/v1"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/shared/configurations/infrastructure"
)

func ConfigRabbitMQ(builder rabbitmqConfigurations.RabbitMQConfigurationBuilder, infra *infrastructure.InfrastructureConfiguration) {
	//add custom message type mappings
	//utils.RegisterCustomMessageTypesToRegistrty(map[string]types.IMessage{"orderCreatedV1": &OrderCreatedV1{}})

	infra.RabbitMQConfigurationBuilder.AddProducer(
		v1.OrderCreatedV1{},
		func(builder producerConfigurations.RabbitMQProducerConfigurationBuilder) {
		})
}
