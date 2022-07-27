package v1

import (
	"context"
	"encoding/json"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/contracts"
	creatingProduct "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/features/creating_product"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/models"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
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
	span.LogFields(log.String("ProductId", command.ProductID))
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
		return nil, err
	}

	response := &creatingProduct.CreateProductResponseDto{ProductID: createdProduct.ProductID}
	bytes, _ := json.Marshal(response)

	c.redisRepository.PutProduct(ctx, createdProduct.ProductID, createdProduct)

	span.LogFields(log.String("CreateProductResponseDto", string(bytes)))

	c.log.Infof("(product created) id: {%s}", command.ProductID)

	return response, nil
}
