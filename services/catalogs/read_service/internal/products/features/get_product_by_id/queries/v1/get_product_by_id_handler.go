package v1

import (
	"context"
	"fmt"
	customErrors "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/http_errors/custom_errors"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/mapper"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/tracing"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/contracts"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/dto"
	getProductByIdDtos "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/features/get_product_by_id/dtos"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/models"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
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
	span.LogFields(log.String("ProductId", query.ProductID.String()))
	span.LogFields(log.Object("Query", query))
	defer span.Finish()

	var product *models.Product
	redisProduct, err := q.redisRepository.GetProduct(ctx, query.ProductID.String())
	if err != nil {
		return nil, tracing.TraceWithErr(span, customErrors.NewApplicationErrorWrap(err, fmt.Sprintf("[GetProductByIdQueryHandler_Handle.GetProduct] error in getting product with id %d in the redis repository", query.ProductID)))
	}

	mongoProduct, err := q.mongoRepository.GetProductById(ctx, query.ProductID)
	if err != nil {
		return nil, tracing.TraceWithErr(span, customErrors.NewApplicationErrorWrap(err, fmt.Sprintf("[GetProductByIdQueryHandler_Handle.GetProduct] error in getting product with id %d in the mongo repository", query.ProductID)))
	}

	if redisProduct != nil {
		product = redisProduct
	} else {
		product = mongoProduct
		err := q.redisRepository.PutProduct(ctx, product.ProductID, product)
		if err != nil {
			return nil, err
		}
	}

	productDto, err := mapper.Map[*dto.ProductDto](product)
	if err != nil {
		return nil, tracing.TraceWithErr(span, customErrors.NewApplicationErrorWrap(err, "[GetProductByIdQueryHandler_Handle.Map] error in the mapping product"))
	}

	q.log.Infow(fmt.Sprintf("[GetProductByIdQueryHandler.Handle] product with id: {%s} fetched", query.ProductID), logger.Fields{"productId": query.ProductID})

	return &getProductByIdDtos.GetProductByIdResponseDto{Product: productDto}, nil
}
