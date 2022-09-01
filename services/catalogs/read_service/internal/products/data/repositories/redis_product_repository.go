package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/tracing"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/models"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

const (
	redisProductPrefixKey = "product_read_service"
)

type redisProductRepository struct {
	log         logger.Logger
	cfg         *config.Config
	redisClient redis.UniversalClient
}

func NewRedisRepository(log logger.Logger, cfg *config.Config, redisClient redis.UniversalClient) *redisProductRepository {
	return &redisProductRepository{log: log, cfg: cfg, redisClient: redisClient}
}

func (r *redisProductRepository) PutProduct(ctx context.Context, key string, product *models.Product) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "redisRepository.PutProduct")
	span.LogFields(log.Object("Model", product))
	span.LogFields(log.String("PrefixKey", r.getRedisProductPrefixKey()))
	defer span.Finish()

	productBytes, err := json.Marshal(product)
	if err != nil {
		return tracing.TraceWithErr(span, errors.Wrap(err, "[redisProductRepository_PutProduct.Marshal] error marshalling product"))
	}

	if err := r.redisClient.HSetNX(ctx, r.getRedisProductPrefixKey(), key, productBytes).Err(); err != nil {
		return tracing.TraceWithErr(span, errors.Wrap(err, fmt.Sprintf("[redisProductRepository_PutProduct.HSetNX] error in updating product with key %s", key)))
	}

	r.log.Infow(fmt.Sprintf("[redisProductRepository.UpdateProduct] result key: %s, prefix: %s", key, r.getRedisProductPrefixKey()), logger.Fields{"Key": key, "PrefixKey": r.getRedisProductPrefixKey()})

	return nil
}

func (r *redisProductRepository) GetProduct(ctx context.Context, key string) (*models.Product, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "redisRepository.GetProduct")
	span.LogFields(log.String("Key", key))
	span.LogFields(log.String("PrefixKey", r.getRedisProductPrefixKey()))
	defer span.Finish()

	productBytes, err := r.redisClient.HGet(ctx, r.getRedisProductPrefixKey(), key).Bytes()
	if err != nil {
		return nil, tracing.TraceWithErr(span, errors.Wrap(err, fmt.Sprintf("[redisProductRepository_GetProduct.HGet] error in getting product with Key %s from database", key)))
	}

	var product models.Product
	if err := json.Unmarshal(productBytes, &product); err != nil {
		return nil, tracing.TraceWithErr(span, err)
	}

	r.log.Infow(fmt.Sprintf("[mongoProductRepository.GetProduct] result: %+v", product), logger.Fields{"Key": key, "PrefixKey": r.getRedisProductPrefixKey()})
	return &product, nil
}

func (r *redisProductRepository) DeleteProduct(ctx context.Context, key string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "redisRepository.GetProduct")
	span.LogFields(log.String("Key", key))
	span.LogFields(log.String("PrefixKey", r.getRedisProductPrefixKey()))
	defer span.Finish()

	if err := r.redisClient.HDel(ctx, r.getRedisProductPrefixKey(), key).Err(); err != nil {
		return tracing.TraceWithErr(span, errors.Wrap(err, fmt.Sprintf("[redisProductRepository_DeleteProduct.HDel] error in deleting product with key %s", key)))
	}

	r.log.Infow(fmt.Sprintf("[redisProductRepository.DeleteProduct] result key: %s, prefix: %s", key, r.getRedisProductPrefixKey()), logger.Fields{"Key": key, "PrefixKey": r.getRedisProductPrefixKey()})
	return nil
}

func (r *redisProductRepository) DeleteAllProducts(ctx context.Context) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "redisRepository.GetProduct")
	span.LogFields(log.String("PrefixKey", r.getRedisProductPrefixKey()))
	defer span.Finish()

	if err := r.redisClient.Del(ctx, r.getRedisProductPrefixKey()).Err(); err != nil {
		return tracing.TraceWithErr(span, errors.Wrap(err, "[redisProductRepository_DeleteAllProducts.Del] error in deleting all products"))
	}

	r.log.Infow(fmt.Sprintf("[redisProductRepository.DeleteAllProducts] result prefix: %s", r.getRedisProductPrefixKey()), logger.Fields{"PrefixKey": r.getRedisProductPrefixKey()})

	return nil
}

func (r *redisProductRepository) getRedisProductPrefixKey() string {
	return redisProductPrefixKey
}
