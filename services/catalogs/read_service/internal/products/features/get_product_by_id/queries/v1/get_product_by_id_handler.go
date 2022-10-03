package v1

import (
	"context"
	"fmt"
	customErrors "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/http_errors/custom_errors"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/mapper"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/otel/tracing"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/otel/tracing/attribute"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/contracts"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/dto"
	getProductByIdDtos "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/features/get_product_by_id/dtos"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/models"
	attribute2 "go.opentelemetry.io/otel/attribute"
)

type GetProductByIdHandler struct {
	log             logger.Logger
	cfg             *config.Config
	mongoRepository contracts.ProductRepository
	redisRepository contracts.ProductCacheRepository
}

func NewGetProductByIdHandler(log logger.Logger, cfg *config.Config, mongoRepository contracts.ProductRepository, redisRepository contracts.ProductCacheRepository) *GetProductByIdHandler {
	return &GetProductByIdHandler{log: log, cfg: cfg, mongoRepository: mongoRepository, redisRepository: redisRepository}
}

func (q *GetProductByIdHandler) Handle(ctx context.Context, query *GetProductById) (*getProductByIdDtos.GetProductByIdResponseDto, error) {
	ctx, span := tracing.Tracer.Start(ctx, "getProductByIdHandler.Handle")
	span.SetAttributes(attribute.Object("Query", query))
	span.SetAttributes(attribute2.String("Id", query.Id.String()))
	defer span.End()

	var product *models.Product
	redisProduct, err := q.redisRepository.GetProduct(ctx, query.Id.String())
	if err != nil {
		return nil, tracing.TraceErrFromSpan(span, customErrors.NewApplicationErrorWrap(err, fmt.Sprintf("[GetProductByIdHandler_Handle.GetProduct] error in getting product with id %d in the redis repository", query.Id)))
	}

	if redisProduct != nil {
		product = redisProduct
	} else {
		mongoProduct, err := q.mongoRepository.GetProductById(ctx, query.Id)
		if err != nil {
			return nil, tracing.TraceErrFromSpan(span, customErrors.NewApplicationErrorWrap(err, fmt.Sprintf("[GetProductByIdHandler_Handle.GetProduct] error in getting product with id %d in the mongo repository", query.Id)))
		}
		product = mongoProduct
		err = q.redisRepository.PutProduct(ctx, product.ProductId, product)
		if err != nil {
			return nil, tracing.TraceErrFromSpan(span, err)
		}
	}

	productDto, err := mapper.Map[*dto.ProductDto](product)
	if err != nil {
		return nil, tracing.TraceErrFromSpan(span, customErrors.NewApplicationErrorWrap(err, "[GetProductByIdHandler_Handle.Map] error in the mapping product"))
	}

	q.log.Infow(fmt.Sprintf("[GetProductByIdHandler.Handle] product with id: {%s} fetched", query.Id), logger.Fields{"ProductId": product.ProductId, "Id": product.Id})

	return &getProductByIdDtos.GetProductByIdResponseDto{Product: productDto}, nil
}
