package configurations

import (
	"context"
	kafkaClient "github.com/mehdihadeli/store-golang-microservice-sample/pkg/kafka"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/delivery"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/features/creating_product/external/consumers"
	consumers2 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/features/deleting_products/external/consumers"
	consumers3 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/features/updating_products/external/consumers"
	"github.com/segmentio/kafka-go"
	"sync"
)

const (
	PoolSize = 30
)

func (c *productsModuleConfigurator) configKafkaConsumers(ctx context.Context) {
	c.Log.Info("Starting Reader Kafka consumers")

	cg := kafkaClient.NewConsumerGroup(c.Cfg.Kafka.Brokers, c.Cfg.Kafka.GroupID, c.Log)
	go cg.ConsumeTopic(ctx, c.getConsumerGroupTopics(), PoolSize, c.processMessages)
}

func (c *productsModuleConfigurator) processMessages(ctx context.Context, r *kafka.Reader, wg *sync.WaitGroup, workerID int) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		message, err := r.FetchMessage(ctx)
		if err != nil {
			c.Log.Warnf("workerID: %v, err: %v", workerID, err)
			continue
		}

		productConsumersBase := delivery.NewProductConsumersBase(c.InfrastructureConfigurations)
		productConsumersBase.LogProcessMessage(message, workerID)

		switch message.Topic {
		case c.Cfg.KafkaTopics.ProductCreated.TopicName:
			consumers.NewCreateProductConsumer(productConsumersBase).Consume(ctx, r, message)
		case c.Cfg.KafkaTopics.ProductUpdated.TopicName:
			consumers3.NewUpdateProductConsumer(productConsumersBase).Consume(ctx, r, message)
		case c.Cfg.KafkaTopics.ProductDeleted.TopicName:
			consumers2.NewDeleteProductConsumer(productConsumersBase).Consume(ctx, r, message)
		}
	}
}

func (c *productsModuleConfigurator) getConsumerGroupTopics() []string {
	return []string{
		c.Cfg.KafkaTopics.ProductCreated.TopicName,
		c.Cfg.KafkaTopics.ProductUpdated.TopicName,
		c.Cfg.KafkaTopics.ProductDeleted.TopicName,
	}
}
