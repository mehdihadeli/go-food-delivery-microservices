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
	uuid "github.com/satori/go.uuid"
)

type CreateProductHandler struct {
	log             logger.Logger
	cfg             *config.Config
	mongoRepository contracts.ProductRepository
	redisRepository contracts.ProductCacheRepository
}

func NewCreateProductHandler(log logger.Logger, cfg *config.Config, mongoRepository contracts.ProductRepository, redisRepository contracts.ProductCacheRepository) *CreateProductHandler {
	return &CreateProductHandler{log: log, cfg: cfg, mongoRepository: mongoRepository, redisRepository: redisRepository}
}

func (c *CreateProductHandler) Handle(ctx context.Context, command *CreateProduct) (*creatingProduct.CreateProductResponseDto, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CreateProductHandler.Handle")
	span.LogFields(log.String("ProductId", command.ProductId))
	span.LogFields(log.Object("Command", command))
	defer span.Finish()

	product := &models.Product{
		Id:          uuid.NewV4().String(), // we generate id ourselves because auto generate mongo string id column with type _id is not an uuid
		ProductId:   command.ProductId,
		Name:        command.Name,
		Description: command.Description,
		Price:       command.Price,
		CreatedAt:   command.CreatedAt,
	}

	createdProduct, err := c.mongoRepository.CreateProduct(ctx, product)
	if err != nil {
		return nil, tracing.TraceWithErr(span, customErrors.NewApplicationErrorWrap(err, "[CreateProductHandler_Handle.CreateProduct] error in creating product in the mongo repository"))
	}

	err = c.redisRepository.PutProduct(ctx, createdProduct.ProductId, createdProduct)
	if err != nil {
		return nil, tracing.TraceWithErr(span, customErrors.NewApplicationErrorWrap(err, "[CreateProductHandler_Handle.PutProduct] error in creating product in the redis repository"))
	}

	response := &creatingProduct.CreateProductResponseDto{ProductId: createdProduct.ProductId}
	span.LogFields(log.Object("CreateProductResponseDto", response))

	c.log.Infow(fmt.Sprintf("[CreateProductHandler.Handle] product with id: {%s} created", command.ProductId), logger.Fields{"productId": command.ProductId})

	return response, nil
}
