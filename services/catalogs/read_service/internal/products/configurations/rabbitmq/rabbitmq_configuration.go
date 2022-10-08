package rabbitmq

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/consumer"
	rabbitmqConfigurations "github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/configurations"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/consumer/configurations"
	creatingProductIntegration "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/features/creating_product/events/integration/external/v1"
	deletingProductIntegration "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/features/deleting_products/events/integration/external/v1"
	updatingProductIntegration "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/features/updating_products/events/integration/external/v1"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/shared/contracts"
)

func ConfigProductsRabbitMQ(builder rabbitmqConfigurations.RabbitMQConfigurationBuilder, infra contracts.InfrastructureConfigurations) {
	//add custom message type mappings
	//utils.RegisterCustomMessageTypesToRegistrty(map[string]types.IMessage{"productCreatedV1": &creatingProductIntegration.ProductCreatedV1{}})

	builder.
		AddConsumer(
			creatingProductIntegration.ProductCreatedV1{},
			func(builder configurations.RabbitMQConsumerConfigurationBuilder) {
				builder.WithHandlers(func(handlersBuilder consumer.ConsumerHandlerConfigurationBuilder) {
					handlersBuilder.AddHandler(creatingProductIntegration.NewProductCreatedConsumer(infra))
				})
			}).
		AddConsumer(
			deletingProductIntegration.ProductDeletedV1{},
			func(builder configurations.RabbitMQConsumerConfigurationBuilder) {
				builder.WithHandlers(func(handlersBuilder consumer.ConsumerHandlerConfigurationBuilder) {
					handlersBuilder.AddHandler(creatingProductIntegration.NewProductCreatedConsumer(infra))
					deletingProductIntegration.NewProductDeletedConsumer(infra)
				})
			}).
		AddConsumer(
			updatingProductIntegration.ProductUpdatedV1{},
			func(builder configurations.RabbitMQConsumerConfigurationBuilder) {
				builder.WithHandlers(func(handlersBuilder consumer.ConsumerHandlerConfigurationBuilder) {
					handlersBuilder.AddHandler(creatingProductIntegration.NewProductCreatedConsumer(infra))
					updatingProductIntegration.NewProductUpdatedConsumer(infra)
				})
			})
}
