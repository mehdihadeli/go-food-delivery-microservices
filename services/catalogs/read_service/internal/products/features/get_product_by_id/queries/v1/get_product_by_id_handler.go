package v1

import (
	"context"
	"fmt"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/mapper"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/contracts"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/dto"
	getProductByIdDtos "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/features/get_product_by_id/dtos"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/models"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

type GetProductByIdQueryHandler struct {
	log             logger.Logger
	cfg             *config.Config
	mongoRepository contracts.ProductRepository
	redisRepository contracts.ProductCacheRepository
}

func NewGetProductByIdQueryHandler(log logger.Logger, cfg *config.Config, mongoRepository contracts.ProductRepository, redisRepository contracts.ProductCacheRepository) *GetProductByIdQueryHandler {
	return &GetProductByIdQueryHandler{log: log, cfg: cfg, mongoRepository: mongoRepository, redisRepository: redisRepository}
}

func (q *GetProductByIdQueryHandler) Handle(ctx context.Context, query *GetProductByIdQuery) (*getProductByIdDtos.GetProductByIdResponseDto, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "getProductByIdHandler.Handle")
	defer span.Finish()

	var product *models.Product
	reidsProduct, err := q.redisRepository.GetProduct(ctx, query.ProductID.String())

	mongoProduct, err := q.mongoRepository.GetProductById(ctx, query.ProductID)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("product with id %s not found", query.ProductID))
	}

	if reidsProduct != nil {
		product = reidsProduct
	} else {
		product = mongoProduct
		err := q.redisRepository.PutProduct(ctx, product.ProductID, product)
		if err != nil {
			return nil, err
		}
	}

	productDto, err := mapper.Map[*dto.ProductDto](product)
	if err != nil {
		return nil, err
	}

	return &getProductByIdDtos.GetProductByIdResponseDto{Product: productDto}, nil
}
