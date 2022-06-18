package updating_product

import (
	"context"
	"github.com/eyazici90/go-mediator/mediator"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/errors"
	kafkaClient "github.com/mehdihadeli/store-golang-microservice-sample/pkg/kafka"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/tracing"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/products/infrastructure/repositories"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/products/mappers"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/products/models"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/products/proto/product_kafka_messages"
	"github.com/opentracing/opentracing-go"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
	"time"
)

//type UpdateProductCmdHandler interface {
//	Handle(ctx context.Context, command *UpdateProduct) error
//}

type updateProductHandler struct {
	log           logger.Logger
	cfg           *config.Config
	pgRepo        repositories.ProductRepository
	kafkaProducer kafkaClient.Producer
}

func NewUpdateProductHandler(log logger.Logger, cfg *config.Config, pgRepo repositories.ProductRepository, kafkaProducer kafkaClient.Producer) *updateProductHandler {
	return &updateProductHandler{log: log, cfg: cfg, pgRepo: pgRepo, kafkaProducer: kafkaProducer}
}

func (c *updateProductHandler) Handle(ctx context.Context, msg mediator.Message) error {

	command, ok := msg.(UpdateProduct)
	if err := errors.CheckType(ok); err != nil {
		return err
	}

	span, ctx := opentracing.StartSpanFromContext(ctx, "updateProductHandler.Handle")
	defer span.Finish()

	productDto := &models.Product{ProductID: command.ProductID, Name: command.Name, Description: command.Description, Price: command.Price}

	product, err := c.pgRepo.UpdateProduct(ctx, productDto)
	if err != nil {
		return err
	}

	evt := &product_kafka_messages.ProductUpdated{Product: mappers.ProductToGrpcMessage(product)}
	msgBytes, err := proto.Marshal(evt)
	if err != nil {
		return err
	}

	message := kafka.Message{
		Topic:   c.cfg.KafkaTopics.ProductUpdated.TopicName,
		Value:   msgBytes,
		Time:    time.Now().UTC(),
		Headers: tracing.GetKafkaTracingHeadersFromSpanCtx(span.Context()),
	}

	return c.kafkaProducer.PublishMessage(ctx, message)
}
