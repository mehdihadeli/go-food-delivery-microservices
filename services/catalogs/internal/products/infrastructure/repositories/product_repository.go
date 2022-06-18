package repositories

import (
	"context"
	models "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/products/models"
	uuid "github.com/satori/go.uuid"
)

type ProductRepository interface {
	CreateProduct(ctx context.Context, product *models.Product) (*models.Product, error)
	UpdateProduct(ctx context.Context, product *models.Product) (*models.Product, error)
	DeleteProductByID(ctx context.Context, uuid uuid.UUID) error
	GetProductById(ctx context.Context, uuid uuid.UUID) (*models.Product, error)
}
