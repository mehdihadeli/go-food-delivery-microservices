package deleting_product

import (
	"context"
	kafkaClient "github.com/mehdihadeli/store-golang-microservice-sample/pkg/kafka"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/tracing"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/contracts"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/contracts/grpc/kafka_messages"
	"github.com/opentracing/opentracing-go"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
	"time"
)

type DeleteProductHandler struct {
	log           logger.Logger
	cfg           *config.Config
	pgRepo        contracts.ProductRepository
	kafkaProducer kafkaClient.Producer
}

func NewDeleteProductHandler(log logger.Logger, cfg *config.Config, pgRepo contracts.ProductRepository, kafkaProducer kafkaClient.Producer) *DeleteProductHandler {
	return &DeleteProductHandler{log: log, cfg: cfg, pgRepo: pgRepo, kafkaProducer: kafkaProducer}
}

func (c *DeleteProductHandler) Handle(ctx context.Context, command DeleteProduct) error {

	span, ctx := opentracing.StartSpanFromContext(ctx, "deleteProductHandler.Handle")
	defer span.Finish()

	if err := c.pgRepo.DeleteProductByID(ctx, command.ProductID); err != nil {
		return err
	}

	evt := &kafka_messages.ProductDeleted{ProductID: command.ProductID.String()}
	msgBytes, err := proto.Marshal(evt)
	if err != nil {
		return err
	}

	message := kafka.Message{
		Topic:   c.cfg.KafkaTopics.ProductDeleted.TopicName,
		Value:   msgBytes,
		Time:    time.Now(),
		Headers: tracing.GetKafkaTracingHeadersFromSpanCtx(span.Context()),
	}

	return c.kafkaProducer.PublishMessage(ctx, message)
}
