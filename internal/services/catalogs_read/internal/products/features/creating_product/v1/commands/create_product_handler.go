package commands

import (
	"context"
	"fmt"

	attribute2 "go.opentelemetry.io/otel/attribute"

	customErrors "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/http_errors/custom_errors"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/otel/tracing"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/otel/tracing/attribute"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/read_service/config"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/read_service/internal/products/contracts"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/read_service/internal/products/features/creating_product/v1/dtos"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/read_service/internal/products/models"
)

type CreateProductHandler struct {
	log             logger.Logger
	cfg             *config.AppConfig
	mongoRepository contracts.ProductRepository
	redisRepository contracts.ProductCacheRepository
}

func NewCreateProductHandler(
	log logger.Logger,
	cfg *config.AppConfig,
	mongoRepository contracts.ProductRepository,
	redisRepository contracts.ProductCacheRepository,
) *CreateProductHandler {
	return &CreateProductHandler{
		log:             log,
		cfg:             cfg,
		mongoRepository: mongoRepository,
		redisRepository: redisRepository,
	}
}

func (c *CreateProductHandler) Handle(
	ctx context.Context,
	command *CreateProduct,
) (*dtos.CreateProductResponseDto, error) {
	ctx, span := tracing.Tracer.Start(ctx, "CreateProductHandler.Handle")
	span.SetAttributes(attribute2.String("ProductId", command.ProductId))
	span.SetAttributes(attribute.Object("Command", command))
	defer span.End()

	product := &models.Product{
		Id:          command.Id, // we generate id ourselves because auto generate mongo string id column with type _id is not an uuid
		ProductId:   command.ProductId,
		Name:        command.Name,
		Description: command.Description,
		Price:       command.Price,
		CreatedAt:   command.CreatedAt,
	}

	createdProduct, err := c.mongoRepository.CreateProduct(ctx, product)
	if err != nil {
		return nil, tracing.TraceErrFromSpan(
			span,
			customErrors.NewApplicationErrorWrap(
				err,
				"[CreateProductHandler_Handle.CreateProduct] error in creating product in the mongo repository",
			),
		)
	}

	err = c.redisRepository.PutProduct(ctx, createdProduct.Id, createdProduct)
	if err != nil {
		return nil, tracing.TraceErrFromSpan(
			span,
			customErrors.NewApplicationErrorWrap(
				err,
				"[CreateProductHandler_Handle.PutProduct] error in creating product in the redis repository",
			),
		)
	}

	response := &dtos.CreateProductResponseDto{Id: createdProduct.Id}
	span.SetAttributes(attribute.Object("CreateProductResponseDto", response))

	c.log.Infow(
		fmt.Sprintf("[CreateProductHandler.Handle] product with id: {%s} created", product.Id),
		logger.Fields{"ProductId": command.ProductId, "Id": product.Id},
	)

	return response, nil
}
