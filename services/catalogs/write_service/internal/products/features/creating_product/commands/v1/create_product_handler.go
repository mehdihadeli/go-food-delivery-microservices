package v1

import (
	"context"
	"fmt"
	customErrors "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/http_errors/custom_errors"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/mapper"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/producer"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/tracing"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/contracts"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/dto"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/creating_product/dtos"
	v1 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/creating_product/events/integration/v1"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/models"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type CreateProductHandler struct {
	log              logger.Logger
	cfg              *config.Config
	repository       contracts.ProductRepository
	rabbitmqProducer producer.Producer
}

func NewCreateProductHandler(log logger.Logger, cfg *config.Config, repository contracts.ProductRepository, rabbitmqProducer producer.Producer) *CreateProductHandler {
	return &CreateProductHandler{log: log, cfg: cfg, repository: repository, rabbitmqProducer: rabbitmqProducer}
}

func (c *CreateProductHandler) Handle(ctx context.Context, command *CreateProduct) (*dtos.CreateProductResponseDto, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CreateProductHandler.Handle")
	span.LogFields(log.String("ProductId", command.ProductID.String()))
	span.LogFields(log.Object("Command", command))
	defer span.Finish()

	product := &models.Product{
		ProductId:   command.ProductID,
		Name:        command.Name,
		Description: command.Description,
		Price:       command.Price,
		CreatedAt:   command.CreatedAt,
	}

	createdProduct, err := c.repository.CreateProduct(ctx, product)
	if err != nil {
		return nil, tracing.TraceWithErr(span, customErrors.NewApplicationErrorWrap(err, "[CreateProductHandler.CreateProduct] error in creating product in the repository"))
	}

	productDto, err := mapper.Map[*dto.ProductDto](createdProduct)
	if err != nil {
		return nil, tracing.TraceWithErr(span, customErrors.NewApplicationErrorWrap(err, "[CreateProductHandler.Map] error in the mapping ProductDto"))
	}

	productCreated := v1.NewProductCreatedV1(productDto)

	err = c.rabbitmqProducer.Publish(ctx, productCreated, nil)
	if err != nil {
		return nil, tracing.TraceWithErr(span, customErrors.NewApplicationErrorWrap(err, "[CreateProductHandler.PublishMessage] error in publishing ProductCreated integration event"))
	}

	c.log.Infow(fmt.Sprintf("[CreateProductHandler.Handle] ProductCreated message with messageId `%s` published to the rabbitmq broker", productCreated.MessageId), logger.Fields{"MessageId": productCreated.MessageId})

	response := &dtos.CreateProductResponseDto{ProductID: product.ProductId}

	span.LogFields(log.Object("CreateProductResponseDto", response))

	c.log.Infow(fmt.Sprintf("[CreateProductHandler.Handle] product with id '%s' created", command.ProductID), logger.Fields{"ProductId": command.ProductID, "MessageId": productCreated.MessageId})

	return response, nil
}
