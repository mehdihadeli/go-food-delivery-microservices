package v1

import (
	"context"
	"github.com/mehdihadeli/go-mediatr"
	kafkaClient "github.com/mehdihadeli/store-golang-microservice-sample/pkg/kafka"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/tracing"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/contracts"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/contracts/proto/kafka_messages"
	"github.com/opentracing/opentracing-go"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
	"time"
)

type DeleteProductCommandHandler struct {
	log           logger.Logger
	cfg           *config.Config
	pgRepo        contracts.ProductRepository
	kafkaProducer kafkaClient.Producer
}

func NewDeleteProductCommandHandler(log logger.Logger, cfg *config.Config, pgRepo contracts.ProductRepository, kafkaProducer kafkaClient.Producer) *DeleteProductCommandHandler {
	return &DeleteProductCommandHandler{log: log, cfg: cfg, pgRepo: pgRepo, kafkaProducer: kafkaProducer}
}

func (c *DeleteProductCommandHandler) Handle(ctx context.Context, command *DeleteProductCommand) (*mediatr.Unit, error) {

	span, ctx := opentracing.StartSpanFromContext(ctx, "deleteProductHandler.Handle")
	defer span.Finish()

	if err := c.pgRepo.DeleteProductByID(ctx, command.ProductID); err != nil {
		return nil, err
	}

	evt := &kafka_messages.ProductDeleted{ProductID: command.ProductID.String()}
	msgBytes, err := proto.Marshal(evt)
	if err != nil {
		return nil, err
	}

	message := kafka.Message{
		Topic:   c.cfg.KafkaTopics.ProductDeleted.TopicName,
		Value:   msgBytes,
		Time:    time.Now(),
		Headers: tracing.GetKafkaTracingHeadersFromSpanCtx(span.Context()),
	}

	return &mediatr.Unit{}, c.kafkaProducer.PublishMessage(ctx, message)
}
