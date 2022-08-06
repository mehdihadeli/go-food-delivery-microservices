package v1

import (
	"context"
	"encoding/json"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/mapper"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/contracts/proto/kafka_messages"
	"time"

	kafkaClient "github.com/mehdihadeli/store-golang-microservice-sample/pkg/kafka"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/tracing"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/contracts"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/creating_product/dtos"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/models"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
)

type CreateProductCommandHandler struct {
	log           logger.Logger
	cfg           *config.Config
	repository    contracts.ProductRepository
	kafkaProducer kafkaClient.Producer
}

func NewCreateProductCommandHandler(log logger.Logger, cfg *config.Config, repository contracts.ProductRepository, kafkaProducer kafkaClient.Producer) *CreateProductCommandHandler {
	return &CreateProductCommandHandler{log: log, cfg: cfg, repository: repository, kafkaProducer: kafkaProducer}
}

func (c *CreateProductCommandHandler) Handle(ctx context.Context, command *CreateProductCommand) (*dtos.CreateProductResponseDto, error) {

	span, ctx := opentracing.StartSpanFromContext(ctx, "CreateProductCommandHandler.Handle")
	span.LogFields(log.String("ProductId", command.ProductID.String()))
	defer span.Finish()

	product := &models.Product{
		ProductID:   command.ProductID,
		Name:        command.Name,
		Description: command.Description,
		Price:       command.Price,
		CreatedAt:   command.CreatedAt,
	}

	createdProduct, err := c.repository.CreateProduct(ctx, product)
	if err != nil {
		return nil, err
	}

	kafkaProduct, err := mapper.Map[*kafka_messages.Product](createdProduct)
	if err != nil {
		return nil, err
	}

	evt := &kafka_messages.ProductCreated{Product: kafkaProduct}
	msgBytes, err := proto.Marshal(evt)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	message := kafka.Message{
		Topic:   c.cfg.KafkaTopics.ProductCreated.TopicName,
		Value:   msgBytes,
		Time:    time.Now(),
		Headers: tracing.GetKafkaTracingHeadersFromSpanCtx(span.Context()),
	}

	err = c.kafkaProducer.PublishMessage(ctx, message)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	response := &dtos.CreateProductResponseDto{ProductID: product.ProductID}
	bytes, _ := json.Marshal(response)

	span.LogFields(log.String("CreateProductResponseDto", string(bytes)))

	return response, nil
}
