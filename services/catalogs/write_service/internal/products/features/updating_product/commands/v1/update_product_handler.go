package v1

import (
	"context"
	"fmt"
	"github.com/mehdihadeli/go-mediatr"
	customErrors "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/http_errors/custom_errors"
	kafkaClient "github.com/mehdihadeli/store-golang-microservice-sample/pkg/kafka"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/mapper"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/tracing"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/contracts"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/contracts/proto/kafka_messages"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/models"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
	"time"
)

type UpdateProductHandler struct {
	log           logger.Logger
	cfg           *config.Config
	pgRepo        contracts.ProductRepository
	kafkaProducer kafkaClient.Producer
}

func NewUpdateProductHandler(log logger.Logger, cfg *config.Config, pgRepo contracts.ProductRepository, kafkaProducer kafkaClient.Producer) *UpdateProductHandler {
	return &UpdateProductHandler{log: log, cfg: cfg, pgRepo: pgRepo, kafkaProducer: kafkaProducer}
}

func (c *UpdateProductHandler) Handle(ctx context.Context, command *UpdateProduct) (*mediatr.Unit, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UpdateProductHandler.Handle")
	span.LogFields(log.String("ProductId", command.ProductID.String()))
	span.LogFields(log.Object("Command", command))
	defer span.Finish()

	p, err := c.pgRepo.GetProductById(ctx, command.ProductID)
	if err != nil {
		return nil, customErrors.NewApplicationErrorWrap(err, fmt.Sprintf("[UpdateProductHandler_Handle.GetProductById] error in fetching product with id %s", command.ProductID))
	}

	if p == nil {
		return nil, customErrors.NewNotFoundErrorWrap(err, fmt.Sprintf("[UpdateProductHandler_Handle.GetProductById] product with id %s not found", command.ProductID))
	}

	product := &models.Product{ProductID: command.ProductID, Name: command.Name, Description: command.Description, Price: command.Price, UpdatedAt: command.UpdatedAt}

	updatedProduct, err := c.pgRepo.UpdateProduct(ctx, product)
	if err != nil {
		return nil, tracing.TraceWithErr(span, customErrors.NewApplicationErrorWrap(err, "[UpdateProductHandler_Handle.UpdateProduct] error in updating product in the repository"))
	}

	productKafka, err := mapper.Map[*kafka_messages.Product](updatedProduct)
	if err != nil {
		return nil, tracing.TraceWithErr(span, customErrors.NewApplicationErrorWrap(err, "[UpdateProductHandler_Handle.Map] error in the mapping product"))
	}

	evt := &kafka_messages.ProductUpdated{Product: productKafka}
	msgBytes, err := proto.Marshal(evt)
	if err != nil {
		return nil, tracing.TraceWithErr(span, customErrors.NewMarshalingErrorWrap(err, "[UpdateProductHandler_Handle.Marshal] error marshalling"))
	}

	message := kafka.Message{
		Topic:   c.cfg.KafkaTopics.ProductUpdated.TopicName,
		Value:   msgBytes,
		Time:    time.Now(),
		Headers: tracing.GetKafkaTracingHeadersFromSpanCtx(span.Context()),
	}

	err = c.kafkaProducer.PublishMessage(ctx, message)
	if err != nil {
		return nil, tracing.TraceWithErr(span, customErrors.NewApplicationErrorWrap(err, "[UpdateProductHandler_Handle.PublishMessage] error in publishing kafka message"))
	}

	c.log.Infow(fmt.Sprintf("[UpdateProductHandler.Handle] product with id: {%s} published to the kafka", command.ProductID), logger.Fields{"productId": command.ProductID})

	c.log.Infow(fmt.Sprintf("[UpdateProductHandler.Handle] product with id: {%s} updated", command.ProductID), logger.Fields{"productId": command.ProductID})

	return &mediatr.Unit{}, nil
}
