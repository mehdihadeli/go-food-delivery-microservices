package configurations

import (
	"context"
	kafkaClient "github.com/mehdihadeli/store-golang-microservice-sample/pkg/kafka"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/delivery"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/features/creating_product"
	"github.com/segmentio/kafka-go"
	"sync"
)

const (
	PoolSize = 30
)

func (pm *ProductModule) configKafkaConsumers(ctx context.Context) {
	pm.Log.Info("Starting Reader Kafka consumers")

	cg := kafkaClient.NewConsumerGroup(pm.Cfg.Kafka.Brokers, pm.Cfg.Kafka.GroupID, pm.Log)
	go cg.ConsumeTopic(ctx, pm.getConsumerGroupTopics(), PoolSize, pm.processMessages)
}

func (pm *ProductModule) processMessages(ctx context.Context, r *kafka.Reader, wg *sync.WaitGroup, workerID int) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		message, err := r.FetchMessage(ctx)
		if err != nil {
			pm.Log.Warnf("workerID: %v, err: %v", workerID, err)
			continue
		}

		productConsumersBase := delivery.NewProductConsumersBase(pm.Infrastructure, pm.Mediator)
		productConsumersBase.LogProcessMessage(message, workerID)

		switch message.Topic {
		case pm.Cfg.KafkaTopics.ProductCreated.TopicName:
			creating_product.NewCreateProductConsumer(productConsumersBase).Consume(ctx, r, message)
			//	s.processProductCreated(ctx, r, m)
			//case s.cfg.KafkaTopics.ProductUpdated.TopicName:
			//	s.processProductUpdated(ctx, r, m)
			//case s.cfg.KafkaTopics.ProductDeleted.TopicName:
			//	s.processProductDeleted(ctx, r, m)
		}
	}
}

func (pm *ProductModule) getConsumerGroupTopics() []string {
	return []string{
		pm.Cfg.KafkaTopics.ProductCreated.TopicName,
		pm.Cfg.KafkaTopics.ProductUpdated.TopicName,
		pm.Cfg.KafkaTopics.ProductDeleted.TopicName,
	}
}
