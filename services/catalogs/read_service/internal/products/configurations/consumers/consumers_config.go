package consumers

import (
	rabbitmqConsumer "github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/consumer"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/consumer/options"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/delivery"
	creatingProductConsumer "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/features/creating_product/consumers"
	creatingProductIntegration "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/features/creating_product/events/integration/external"
	deletingProductConsumer "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/features/deleting_products/consumers"
	deletingProductIntegration "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/features/deleting_products/events/integration/external"
	updatingProductConsumer "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/features/updating_products/consumers"
	updatingProductIntegration "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/features/updating_products/events/integration/external"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/shared/configurations/infrastructure"
)

func ConfigConsumers(infra *infrastructure.InfrastructureConfigurations) error {
	consumerBase := delivery.NewProductConsumersBase(infra)

	productCreatedConsumer, err := rabbitmqConsumer.NewRabbitMQConsumer[*creatingProductIntegration.ProductCreated](
		infra.RabbitMQConnection,
		func(builder *options.RabbitMQConsumerOptionsBuilder[*creatingProductIntegration.ProductCreated]) {},
		infra.EventSerializer,
		infra.Log,
		creatingProductConsumer.NewProductCreatedConsumer(consumerBase))
	if err != nil {
		return err
	}
	infra.Consumers = append(infra.Consumers, productCreatedConsumer)

	productDeletedConsumer, err := rabbitmqConsumer.NewRabbitMQConsumer[*deletingProductIntegration.ProductDeleted](
		infra.RabbitMQConnection,
		func(builder *options.RabbitMQConsumerOptionsBuilder[*deletingProductIntegration.ProductDeleted]) {},
		infra.EventSerializer,
		infra.Log,
		deletingProductConsumer.NewProductDeletedConsumer(consumerBase))
	if err != nil {
		return err
	}
	infra.Consumers = append(infra.Consumers, productDeletedConsumer)

	productUpdatedConsumer, err := rabbitmqConsumer.NewRabbitMQConsumer[*updatingProductIntegration.ProductUpdated](
		infra.RabbitMQConnection,
		func(builder *options.RabbitMQConsumerOptionsBuilder[*updatingProductIntegration.ProductUpdated]) {},
		infra.EventSerializer,
		infra.Log,
		updatingProductConsumer.NewProductUpdatedConsumer(consumerBase))
	if err != nil {
		return err
	}
	infra.Consumers = append(infra.Consumers, productUpdatedConsumer)

	return nil
}
