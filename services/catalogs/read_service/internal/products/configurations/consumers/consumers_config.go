package consumers

import (
	rabbitmqConsumer "github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/consumer"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/consumer/options"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/delivery"
	creatingProductIntegration "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/features/creating_product/events/integration/external/v1"
	deletingProductIntegration "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/features/deleting_products/events/integration/external/v1"
	updatingProductIntegration "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/features/updating_products/events/integration/external/v1"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/shared/configurations/infrastructure"
)

func ConfigConsumers(infra *infrastructure.InfrastructureConfigurations) error {
	consumerBase := delivery.NewProductConsumersBase(infra)

	//add custom message type mappings
	//utils.RegisterCustomMessageTypesToRegistrty(map[string]types.IMessage{"productCreatedV1": &creatingProductIntegration.ProductCreatedV1{}})

	productCreatedConsumer, err := rabbitmqConsumer.NewRabbitMQConsumer[*creatingProductIntegration.ProductCreatedV1](
		infra.RabbitMQConnection,
		func(builder *options.RabbitMQConsumerOptionsBuilder[*creatingProductIntegration.ProductCreatedV1]) {},
		infra.EventSerializer,
		infra.Log,
		creatingProductIntegration.NewProductCreatedConsumer(consumerBase))
	if err != nil {
		return err
	}
	infra.Consumers = append(infra.Consumers, productCreatedConsumer)

	productDeletedConsumer, err := rabbitmqConsumer.NewRabbitMQConsumer[*deletingProductIntegration.ProductDeletedV1](
		infra.RabbitMQConnection,
		func(builder *options.RabbitMQConsumerOptionsBuilder[*deletingProductIntegration.ProductDeletedV1]) {},
		infra.EventSerializer,
		infra.Log,
		deletingProductIntegration.NewProductDeletedConsumer(consumerBase))
	if err != nil {
		return err
	}
	infra.Consumers = append(infra.Consumers, productDeletedConsumer)

	productUpdatedConsumer, err := rabbitmqConsumer.NewRabbitMQConsumer[*updatingProductIntegration.ProductUpdatedV1](
		infra.RabbitMQConnection,
		func(builder *options.RabbitMQConsumerOptionsBuilder[*updatingProductIntegration.ProductUpdatedV1]) {},
		infra.EventSerializer,
		infra.Log,
		updatingProductIntegration.NewProductUpdatedConsumer(consumerBase))
	if err != nil {
		return err
	}
	infra.Consumers = append(infra.Consumers, productUpdatedConsumer)

	return nil
}
