package contracts

import (
	"context"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/models"
)

type ProductCacheRepository interface {
	PutProduct(ctx context.Context, key string, product *models.Product) error
	GetProductById(ctx context.Context, key string) (*models.Product, error)
	DeleteProduct(ctx context.Context, key string) error
	DeleteAllProducts(ctx context.Context) error
}
