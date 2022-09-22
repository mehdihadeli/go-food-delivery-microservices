package repositories

import (
	"context"
	"emperror.dev/errors"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/tracing"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/contracts"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/models"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

const (
	redisProductPrefixKey = "product_read_service"
)

type redisProductRepository struct {
	log         logger.Logger
	cfg         *config.Config
	redisClient redis.UniversalClient
}

func NewRedisRepository(log logger.Logger, cfg *config.Config, redisClient redis.UniversalClient) contracts.ProductCacheRepository {
	return &redisProductRepository{log: log, cfg: cfg, redisClient: redisClient}
}

func (r *redisProductRepository) PutProduct(ctx context.Context, key string, product *models.Product) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "redisRepository.PutProduct")
	span.LogFields(log.String("PrefixKey", r.getRedisProductPrefixKey()))
	span.LogFields(log.String("Key", key))
	defer span.Finish()

	productBytes, err := json.Marshal(product)
	if err != nil {
		return tracing.TraceWithErr(span, errors.WrapIf(err, "[redisProductRepository_PutProduct.Marshal] error marshalling product"))
	}

	if err := r.redisClient.HSetNX(ctx, r.getRedisProductPrefixKey(), key, productBytes).Err(); err != nil {
		return tracing.TraceWithErr(span, errors.WrapIf(err, fmt.Sprintf("[redisProductRepository_PutProduct.HSetNX] error in updating product with key %s", key)))
	}

	span.LogFields(log.Object("Product", product))
	r.log.Infow(fmt.Sprintf("[redisProductRepository.PutProduct] product with key '%s', prefix '%s'  updated successfully",
		key,
		r.getRedisProductPrefixKey()),
		logger.Fields{"Product": product, "Id": product.ProductId, "Key": key, "PrefixKey": r.getRedisProductPrefixKey()})

	return nil
}

func (r *redisProductRepository) GetProduct(ctx context.Context, key string) (*models.Product, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "redisRepository.GetProduct")
	span.LogFields(log.String("Key", key))
	span.LogFields(log.String("PrefixKey", r.getRedisProductPrefixKey()))
	defer span.Finish()

	productBytes, err := r.redisClient.HGet(ctx, r.getRedisProductPrefixKey(), key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}

		return nil, tracing.TraceWithErr(span, errors.WrapIf(err, fmt.Sprintf("[redisProductRepository_GetProduct.HGet] error in getting product with Key %s from database", key)))
	}

	var product models.Product
	if err := json.Unmarshal(productBytes, &product); err != nil {
		return nil, tracing.TraceWithErr(span, err)
	}

	span.LogFields(log.Object("Product", product))
	r.log.Infow(fmt.Sprintf("[redisProductRepository.GetProduct] product with with key '%s', prefix '%s' laoded", key, r.getRedisProductPrefixKey()),
		logger.Fields{"Product": product, "Id": product.ProductId, "Key": key, "PrefixKey": r.getRedisProductPrefixKey()})

	return &product, nil
}

func (r *redisProductRepository) DeleteProduct(ctx context.Context, key string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "redisRepository.GetProduct")
	span.LogFields(log.String("Key", key))
	span.LogFields(log.String("PrefixKey", r.getRedisProductPrefixKey()))
	defer span.Finish()

	if err := r.redisClient.HDel(ctx, r.getRedisProductPrefixKey(), key).Err(); err != nil {
		return tracing.TraceWithErr(span, errors.WrapIf(err, fmt.Sprintf("[redisProductRepository_DeleteProduct.HDel] error in deleting product with key %s", key)))
	}

	r.log.Infow(fmt.Sprintf("[redisProductRepository.DeleteProduct] product with key %s, prefix: %s deleted successfully", key, r.getRedisProductPrefixKey()), logger.Fields{"Key": key, "PrefixKey": r.getRedisProductPrefixKey()})

	return nil
}

func (r *redisProductRepository) DeleteAllProducts(ctx context.Context) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "redisRepository.GetProduct")
	span.LogFields(log.String("PrefixKey", r.getRedisProductPrefixKey()))
	defer span.Finish()

	if err := r.redisClient.Del(ctx, r.getRedisProductPrefixKey()).Err(); err != nil {
		return tracing.TraceWithErr(span, errors.WrapIf(err, "[redisProductRepository_DeleteAllProducts.Del] error in deleting all products"))
	}

	r.log.Infow("[redisProductRepository.DeleteAllProducts] all products deleted", logger.Fields{"PrefixKey": r.getRedisProductPrefixKey()})

	return nil
}

func (r *redisProductRepository) getRedisProductPrefixKey() string {
	return redisProductPrefixKey
}
