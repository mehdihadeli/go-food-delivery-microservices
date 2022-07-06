package configurations

import (
	"context"
	kafkaClient "github.com/mehdihadeli/store-golang-microservice-sample/pkg/kafka"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/mediatr"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/delivery"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/features/creating_product"
	"github.com/segmentio/kafka-go"
	"sync"
)

const (
	PoolSize = 30
)

func (c *productsModuleConfigurator) configKafkaConsumers(ctx context.Context, mediator *mediatr.Mediator) {
	c.Log.Info("Starting Reader Kafka consumers")

	cg := kafkaClient.NewConsumerGroup(c.Cfg.Kafka.Brokers, c.Cfg.Kafka.GroupID, c.Log)
	go cg.ConsumeTopic(ctx, mediator, c.getConsumerGroupTopics(), PoolSize, c.processMessages)
}

func (c *productsModuleConfigurator) processMessages(ctx context.Context, mediator *mediatr.Mediator, r *kafka.Reader, wg *sync.WaitGroup, workerID int) {
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

		productConsumersBase := delivery.NewProductConsumersBase(c.InfrastructureConfigurations, mediator)
		productConsumersBase.LogProcessMessage(message, workerID)

		switch message.Topic {
		case c.Cfg.KafkaTopics.ProductCreated.TopicName:
			creating_product.NewCreateProductConsumer(productConsumersBase).Consume(ctx, r, message)
			//	s.processProductCreated(ctx, r, m)
			//case s.cfg.KafkaTopics.ProductUpdated.TopicName:
			//	s.processProductUpdated(ctx, r, m)
			//case s.cfg.KafkaTopics.ProductDeleted.TopicName:
			//	s.processProductDeleted(ctx, r, m)
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
