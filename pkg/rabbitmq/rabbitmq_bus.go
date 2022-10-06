package rabbitmq

import (
	"context"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/core/serializer"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/bus"
	messageConsumer "github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/consumer"
	rabbitmqBus "github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/bus"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/configurations"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/consumer"
	rabbitmqProducer "github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/producer"
	producerConfigurations "github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/producer/configurations"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/types"
	"github.com/solsw/go2linq/v2"
)

func AddRabbitMQBus(ctx context.Context, cfg *config.RabbitMQConfig, rabbitmqBuilderFunc configurations.RabbitMQConfigurationBuilderFuc, serializer serializer.EventSerializer, log logger.Logger) (bus.Bus, error) {
	builder := configurations.NewRabbitMQConfigurationBuilder()
	rabbitmqBuilderFunc(builder)

	rabbitmqConfiguration := builder.Build()

	conn, err := types.NewRabbitMQConnection(ctx, cfg)
	if err != nil {
		return nil, err
	}

	producers := go2linq.ToMapMust(go2linq.NewEnSlice(rabbitmqConfiguration.ProducersConfigurations...), func(source *producerConfigurations.RabbitMQProducerConfiguration) string {
		return source.ProducerMessageType.String()
	})

	p, err := rabbitmqProducer.NewRabbitMQProducer(conn, producers, log, serializer)
	if err != nil {
		return nil, err
	}

	var consumers []messageConsumer.Consumer
	for _, consumerConfiguration := range rabbitmqConfiguration.ConsumersConfigurations {
		consumer, err := consumer.NewRabbitMQConsumer(consumerConfiguration.ConsumerMessageType, serializer, log, conn, consumerConfiguration)
		if err != nil {
			return nil, err
		}
		consumers = append(consumers, consumer)
	}

	b := rabbitmqBus.NewRabbitMQBus(log, p, consumers)

	return b, nil
}
