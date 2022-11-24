package rabbitmq

import (
	rabbitmqConfigurations "github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/configurations"
	producerConfigurations "github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/producer/configurations"
	createOrderIntegrationEventsV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/features/creating_order/v1/events/integration_events"
)

func ConfigOrdersRabbitMQ(builder rabbitmqConfigurations.RabbitMQConfigurationBuilder) {
	//add custom message type mappings
	//utils.RegisterCustomMessageTypesToRegistrty(map[string]types.IMessage{"orderCreatedV1": &OrderCreatedV1{}})

	builder.AddProducer(
		createOrderIntegrationEventsV1.OrderCreatedV1{},
		func(builder producerConfigurations.RabbitMQProducerConfigurationBuilder) {
		})
}
