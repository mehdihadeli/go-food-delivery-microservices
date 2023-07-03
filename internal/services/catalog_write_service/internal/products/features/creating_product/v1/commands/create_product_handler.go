package createProductCommand

import (
	"context"
	"fmt"
	"net/http"

	attribute2 "go.opentelemetry.io/otel/attribute"

	customErrors "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/http_errors/custom_errors"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/mapper"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/messaging/producer"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/otel/tracing"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/otel/tracing/attribute"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/contracts/data"
	dtoV1 "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/dto/v1"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/features/creating_product/v1/dtos"
	integrationEvents "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/features/creating_product/v1/events/integration_events"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/models"
)

type CreateProductHandler struct {
	log              logger.Logger
	uow              data.CatalogUnitOfWork
	rabbitmqProducer producer.Producer
	tracer           tracing.AppTracer
}

func NewCreateProductHandler(
	log logger.Logger,
	uow data.CatalogUnitOfWork,
	rabbitmqProducer producer.Producer,
	tracer tracing.AppTracer,
) *CreateProductHandler {
	return &CreateProductHandler{
		log:              log,
		uow:              uow,
		rabbitmqProducer: rabbitmqProducer,
		tracer:           tracer,
	}
}

func (c *CreateProductHandler) Handle(
	ctx context.Context,
	command *CreateProduct,
) (*dtos.CreateProductResponseDto, error) {
	ctx, span := c.tracer.Start(ctx, "CreateProductHandler.Handle")
	span.SetAttributes(attribute2.String("ProductId", command.ProductID.String()))
	span.SetAttributes(attribute.Object("Command", command))
	defer span.End()

	product := &models.Product{
		ProductId:   command.ProductID,
		Name:        command.Name,
		Description: command.Description,
		Price:       command.Price,
		CreatedAt:   command.CreatedAt,
	}

	var createProductResult *dtos.CreateProductResponseDto

	err := c.uow.Do(ctx, func(catalogContext data.CatalogContext) error {
		createdProduct, err := catalogContext.Products().CreateProduct(ctx, product)
		if err != nil {
			return tracing.TraceErrFromSpan(
				span,
				customErrors.NewApplicationErrorWrapWithCode(
					err,
					http.StatusConflict,
					"[CreateProductHandler.CreateProduct] product already exists",
				),
			)
		}
		productDto, err := mapper.Map[*dtoV1.ProductDto](createdProduct)
		if err != nil {
			return tracing.TraceErrFromSpan(
				span,
				customErrors.NewApplicationErrorWrap(
					err,
					"[CreateProductHandler.Map] error in the mapping ProductDto",
				),
			)
		}

		productCreated := integrationEvents.NewProductCreatedV1(productDto)

		err = c.rabbitmqProducer.PublishMessage(ctx, productCreated, nil)
		if err != nil {
			return tracing.TraceErrFromSpan(
				span,
				customErrors.NewApplicationErrorWrap(
					err,
					"[CreateProductHandler.PublishMessage] error in publishing ProductCreated integration_events event",
				),
			)
		}

		c.log.Infow(
			fmt.Sprintf(
				"[CreateProductHandler.Handle] ProductCreated message with messageId `%s` published to the rabbitmq broker",
				productCreated.MessageId,
			),
			logger.Fields{"MessageId": productCreated.MessageId},
		)

		createProductResult = &dtos.CreateProductResponseDto{ProductID: product.ProductId}

		span.SetAttributes(attribute.Object("CreateProductResultDto", createProductResult))

		c.log.Infow(
			fmt.Sprintf(
				"[CreateProductHandler.Handle] product with id '%s' created",
				command.ProductID,
			),
			logger.Fields{"ProductId": command.ProductID, "MessageId": productCreated.MessageId},
		)

		return nil
	})

	return createProductResult, err
}
