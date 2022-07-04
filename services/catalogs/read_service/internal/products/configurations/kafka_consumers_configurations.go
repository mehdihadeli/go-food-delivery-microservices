package configurations

import (
	"context"
	kafkaClient "github.com/mehdihadeli/store-golang-microservice-sample/pkg/kafka"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/mediatr"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/features/creating_product"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/shared/configurations"
	"github.com/segmentio/kafka-go"
	"sync"
)

const (
	PoolSize = 30
)

type ProductKafkaConsumersConfigurator struct {
	*ProductModuleConfigurations
}

type ProductKafkaConsumersConfigurations struct {
	*configurations.Infrastructure
	*mediatr.Mediator
}

func (pc *ProductKafkaConsumersConfigurator) configKafkaConsumers(ctx context.Context) {

	consumerConfigurations := &ProductKafkaConsumersConfigurations{Infrastructure: pc.Infrastructure, Mediator: pc.Mediator}

	pc.Infrastructure.Log.Info("Starting Reader Kafka consumers")
	cg := kafkaClient.NewConsumerGroup(pc.Infrastructure.Cfg.Kafka.Brokers, pc.Infrastructure.Cfg.Kafka.GroupID, pc.Infrastructure.Log)
	go cg.ConsumeTopic(ctx, consumerConfigurations.getConsumerGroupTopics(), PoolSize, consumerConfigurations.processMessages)
}

func (pm *ProductKafkaConsumersConfigurations) processMessages(ctx context.Context, r *kafka.Reader, wg *sync.WaitGroup, workerID int) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		m, err := r.FetchMessage(ctx)
		if err != nil {
			pm.Infrastructure.Log.Warnf("workerID: %v, err: %v", workerID, err)
			continue
		}

		pm.LogProcessMessage(m, workerID)

		switch m.Topic {
		case pm.Infrastructure.Cfg.KafkaTopics.ProductCreated.TopicName:
			creating_product.NewCreateProductConsumer(pm).Process(ctx, r, m)
			//	s.processProductCreated(ctx, r, m)
			//case s.cfg.KafkaTopics.ProductUpdated.TopicName:
			//	s.processProductUpdated(ctx, r, m)
			//case s.cfg.KafkaTopics.ProductDeleted.TopicName:
			//	s.processProductDeleted(ctx, r, m)
		}
	}
}

func (pm *ProductKafkaConsumersConfigurations) getConsumerGroupTopics() []string {
	return []string{
		pm.Infrastructure.Cfg.KafkaTopics.ProductCreated.TopicName,
		pm.Infrastructure.Cfg.KafkaTopics.ProductUpdated.TopicName,
		pm.Infrastructure.Cfg.KafkaTopics.ProductDeleted.TopicName,
	}
}

func (pm *ProductKafkaConsumersConfigurations) CommitMessage(ctx context.Context, r *kafka.Reader, m kafka.Message) {
	pm.Infrastructure.Metrics.SuccessKafkaMessages.Inc()
	pm.Infrastructure.Log.KafkaLogCommittedMessage(m.Topic, m.Partition, m.Offset)

	if err := r.CommitMessages(ctx, m); err != nil {
		pm.Infrastructure.Log.WarnMsg("commitMessage", err)
	}
}

func (pm *ProductKafkaConsumersConfigurations) LogProcessMessage(m kafka.Message, workerID int) {
	pm.Infrastructure.Log.KafkaProcessMessage(m.Topic, m.Partition, string(m.Value), workerID, m.Offset, m.Time)
}

func (pm *ProductKafkaConsumersConfigurations) CommitErrMessage(ctx context.Context, r *kafka.Reader, m kafka.Message) {
	pm.Infrastructure.Metrics.ErrorKafkaMessages.Inc()
	pm.Infrastructure.Log.KafkaLogCommittedMessage(m.Topic, m.Partition, m.Offset)
	if err := r.CommitMessages(ctx, m); err != nil {
		pm.Infrastructure.Log.WarnMsg("commitMessage", err)
	}
}
