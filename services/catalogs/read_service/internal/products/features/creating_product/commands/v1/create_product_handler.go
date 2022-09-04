package v1

import (
	"context"
	"fmt"
	customErrors "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/http_errors/custom_errors"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/tracing"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/contracts"
	creatingProduct "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/features/creating_product"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/models"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type CreateProductCommandHandler struct {
	log             logger.Logger
	cfg             *config.Config
	mongoRepository contracts.ProductRepository
	redisRepository contracts.ProductCacheRepository
}

func NewCreateProductCommandHandler(log logger.Logger, cfg *config.Config, mongoRepository contracts.ProductRepository, redisRepository contracts.ProductCacheRepository) *CreateProductCommandHandler {
	return &CreateProductCommandHandler{log: log, cfg: cfg, mongoRepository: mongoRepository, redisRepository: redisRepository}
}

func (c *CreateProductCommandHandler) Handle(ctx context.Context, command *CreateProductCommand) (*creatingProduct.CreateProductResponseDto, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CreateProductCommandHandler.Handle")
	span.LogFields(log.String("ProductId", command.ProductID))
	span.LogFields(log.Object("Command", command))
	defer span.Finish()

	product := &models.Product{
		ProductID:   command.ProductID,
		Name:        command.Name,
		Description: command.Description,
		Price:       command.Price,
		CreatedAt:   command.CreatedAt,
	}

	createdProduct, err := c.mongoRepository.CreateProduct(ctx, product)
	if err != nil {
		return nil, tracing.TraceWithErr(span, customErrors.NewApplicationErrorWrap(err, "[CreateProductCommandHandler_Handle.CreateProduct] error in creating product in the mongo repository"))
	}

	err = c.redisRepository.PutProduct(ctx, createdProduct.ProductID, createdProduct)
	if err != nil {
		return nil, tracing.TraceWithErr(span, customErrors.NewApplicationErrorWrap(err, "[CreateProductCommandHandler_Handle.PutProduct] error in creating product in the redis repository"))
	}

	response := &creatingProduct.CreateProductResponseDto{ProductID: createdProduct.ProductID}
	span.LogFields(log.Object("CreateProductResponseDto", response))

	c.log.Infow(fmt.Sprintf("[CreateProductCommandHandler.Handle] product with id: {%s} created", command.ProductID), logger.Fields{"productId": command.ProductID})

	return response, nil
}
