package rabbitmq

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/consumer"
	rabbitmqConfigurations "github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/configurations"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/consumer/configurations"
	createProductExternalEventV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/features/creating_product/v1/events/integration_events/external_events"
	deleteProductExternalEventV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/features/deleting_products/v1/events/integration_events/external_events"
	updateProductExternalEventsV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/features/updating_products/v1/events/integration_events/external_events"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/shared/contracts"
)

func ConfigProductsRabbitMQ(builder rabbitmqConfigurations.RabbitMQConfigurationBuilder, infra *contracts.InfrastructureConfigurations) {
	// add custom message type mappings
	// utils.RegisterCustomMessageTypesToRegistrty(map[string]types.IMessage{"productCreatedV1": &creatingProductIntegration.ProductCreatedV1{}})

	builder.
		AddConsumer(
			createProductExternalEventV1.ProductCreatedV1{},
			func(builder configurations.RabbitMQConsumerConfigurationBuilder) {
				builder.WithHandlers(func(handlersBuilder consumer.ConsumerHandlerConfigurationBuilder) {
					handlersBuilder.AddHandler(createProductExternalEventV1.NewProductCreatedConsumer(infra))
				})
			}).
		AddConsumer(
			deleteProductExternalEventV1.ProductDeletedV1{},
			func(builder configurations.RabbitMQConsumerConfigurationBuilder) {
				builder.WithHandlers(func(handlersBuilder consumer.ConsumerHandlerConfigurationBuilder) {
					handlersBuilder.AddHandler(deleteProductExternalEventV1.NewProductDeletedConsumer(infra))
					deleteProductExternalEventV1.NewProductDeletedConsumer(infra)
				})
			}).
		AddConsumer(
			updateProductExternalEventsV1.ProductUpdatedV1{},
			func(builder configurations.RabbitMQConsumerConfigurationBuilder) {
				builder.WithHandlers(func(handlersBuilder consumer.ConsumerHandlerConfigurationBuilder) {
					handlersBuilder.AddHandler(updateProductExternalEventsV1.NewProductUpdatedConsumer(infra))
					updateProductExternalEventsV1.NewProductUpdatedConsumer(infra)
				})
			})
}
