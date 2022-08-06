package v1

import (
	"context"
	"fmt"
	"github.com/mehdihadeli/go-mediatr"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/http_errors"
	kafkaClient "github.com/mehdihadeli/store-golang-microservice-sample/pkg/kafka"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/mapper"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/tracing"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/contracts"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/contracts/proto/kafka_messages"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/models"
	"github.com/opentracing/opentracing-go"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
	"time"
)

type UpdateProductCommandHandler struct {
	log           logger.Logger
	cfg           *config.Config
	pgRepo        contracts.ProductRepository
	kafkaProducer kafkaClient.Producer
}

func NewUpdateProductCommandHandler(log logger.Logger, cfg *config.Config, pgRepo contracts.ProductRepository, kafkaProducer kafkaClient.Producer) *UpdateProductCommandHandler {
	return &UpdateProductCommandHandler{log: log, cfg: cfg, pgRepo: pgRepo, kafkaProducer: kafkaProducer}
}

func (c *UpdateProductCommandHandler) Handle(ctx context.Context, command *UpdateProductCommand) (*mediatr.Unit, error) {

	span, ctx := opentracing.StartSpanFromContext(ctx, "UpdateProductCommandHandler.Handle")
	defer span.Finish()

	_, err := c.pgRepo.GetProductById(ctx, command.ProductID)

	if err != nil {
		return nil, httpErrors.NewNotFoundError(fmt.Sprintf("product with id %s not found", command.ProductID))
	}

	product := &models.Product{ProductID: command.ProductID, Name: command.Name, Description: command.Description, Price: command.Price, UpdatedAt: command.UpdatedAt}

	updatedProduct, err := c.pgRepo.UpdateProduct(ctx, product)
	if err != nil {
		return nil, err
	}

	productKafka, err := mapper.Map[*kafka_messages.Product](updatedProduct)
	if err != nil {
		return nil, err
	}

	evt := &kafka_messages.ProductUpdated{Product: productKafka}
	msgBytes, err := proto.Marshal(evt)
	if err != nil {
		return nil, err
	}

	message := kafka.Message{
		Topic:   c.cfg.KafkaTopics.ProductUpdated.TopicName,
		Value:   msgBytes,
		Time:    time.Now(),
		Headers: tracing.GetKafkaTracingHeadersFromSpanCtx(span.Context()),
	}

	return &mediatr.Unit{}, c.kafkaProducer.PublishMessage(ctx, message)
}
