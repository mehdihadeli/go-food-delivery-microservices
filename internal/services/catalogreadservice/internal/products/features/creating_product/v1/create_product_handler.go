package v1

import (
	"context"
	"fmt"

	customErrors "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/http/httperrors/customerrors"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/logger"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/otel/tracing"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogreadservice/internal/products/contracts/data"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogreadservice/internal/products/features/creating_product/v1/dtos"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogreadservice/internal/products/models"
)

type CreateProductHandler struct {
	log             logger.Logger
	mongoRepository data.ProductRepository
	redisRepository data.ProductCacheRepository
	tracer          tracing.AppTracer
}

func NewCreateProductHandler(
	log logger.Logger,
	mongoRepository data.ProductRepository,
	redisRepository data.ProductCacheRepository,
	tracer tracing.AppTracer,
) *CreateProductHandler {
	return &CreateProductHandler{
		log:             log,
		mongoRepository: mongoRepository,
		redisRepository: redisRepository,
		tracer:          tracer,
	}
}

func (c *CreateProductHandler) Handle(
	ctx context.Context,
	command *CreateProduct,
) (*dtos.CreateProductResponseDto, error) {
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
		return nil, customErrors.NewApplicationErrorWrap(
			err,
			"error in creating product in the mongo repository",
		)
	}

	err = c.redisRepository.PutProduct(ctx, createdProduct.Id, createdProduct)
	if err != nil {
		return nil, customErrors.NewApplicationErrorWrap(
			err,
			"error in creating product in the redis repository",
		)
	}

	response := &dtos.CreateProductResponseDto{Id: createdProduct.Id}

	c.log.Infow(
		fmt.Sprintf(
			"product with id: {%s} created",
			product.Id,
		),
		logger.Fields{"ProductId": command.ProductId, "Id": product.Id},
	)

	return response, nil
}
