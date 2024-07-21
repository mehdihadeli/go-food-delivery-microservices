package data

import (
	"context"

	"github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogreadservice/internal/products/models"
)

type ProductCacheRepository interface {
	PutProduct(ctx context.Context, key string, product *models.Product) error
	GetProductById(ctx context.Context, key string) (*models.Product, error)
	DeleteProduct(ctx context.Context, key string) error
	DeleteAllProducts(ctx context.Context) error
}
