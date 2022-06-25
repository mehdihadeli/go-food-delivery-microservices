package repositories

import (
	"context"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/products/models"
	uuid "github.com/satori/go.uuid"
)

type ProductRepository interface {
	GetAllProducts(ctx context.Context) ([]*models.Product, error)
	GetProductsByPage(ctx context.Context, page int, skip int) ([]*models.Product, error)
	GetProductById(ctx context.Context, uuid uuid.UUID) (*models.Product, error)
	CreateProduct(ctx context.Context, product *models.Product) (*models.Product, error)
	UpdateProduct(ctx context.Context, product *models.Product) (*models.Product, error)
	DeleteProductByID(ctx context.Context, uuid uuid.UUID) error
}
