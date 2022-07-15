package contracts

import (
	"context"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/models"
)

type ProductCacheRepository interface {
	PutProduct(ctx context.Context, key string, product *models.Product)
	GetProduct(ctx context.Context, key string) (*models.Product, error)
	DelProduct(ctx context.Context, key string)
	DelAllProducts(ctx context.Context)
}
