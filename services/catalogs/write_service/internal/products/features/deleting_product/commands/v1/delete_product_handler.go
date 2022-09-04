package v1

import (
	"context"
	"fmt"
	"github.com/mehdihadeli/go-mediatr"
	customErrors "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/http_errors/custom_errors"
	kafkaClient "github.com/mehdihadeli/store-golang-microservice-sample/pkg/kafka"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/tracing"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/contracts"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/contracts/proto/kafka_messages"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
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
	span.LogFields(log.String("ProductId", command.ProductID.String()))
	span.LogFields(log.Object("Command", command))
	defer span.Finish()

	if err := c.pgRepo.DeleteProductByID(ctx, command.ProductID); err != nil {
		return nil, tracing.TraceWithErr(span, customErrors.NewApplicationErrorWrap(err, "[DeleteProductCommandHandler_Handle.DeleteProductByID] error in deleting product in the repository"))
	}

	evt := &kafka_messages.ProductDeleted{ProductID: command.ProductID.String()}
	msgBytes, err := proto.Marshal(evt)
	if err != nil {
		return nil, tracing.TraceWithErr(span, customErrors.NewApplicationErrorWrap(err, "[DeleteProductCommandHandler_Handle.Marshal] error in marshaling proto event ProductDeleted"))
	}

	message := kafka.Message{
		Topic:   c.cfg.KafkaTopics.ProductDeleted.TopicName,
		Value:   msgBytes,
		Time:    time.Now(),
		Headers: tracing.GetKafkaTracingHeadersFromSpanCtx(span.Context()),
	}

	err = c.kafkaProducer.PublishMessage(ctx, message)
	if err != nil {
		return nil, tracing.TraceWithErr(span, customErrors.NewApplicationErrorWrap(err, "[DeleteProductCommandHandler_Handle.PublishMessage] error in publishing kafka message"))
	}

	c.log.Infow(fmt.Sprintf("[DeleteProductCommandHandler.Handle] product with id: {%s} deleted", command.ProductID), logger.Fields{"productId": command.ProductID})

	return &mediatr.Unit{}, nil
}
