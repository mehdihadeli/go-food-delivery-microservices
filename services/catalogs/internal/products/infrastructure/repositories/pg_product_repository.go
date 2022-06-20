package repositories

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/products/models"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

type postgresProductRepository struct {
	log logger.Logger
	cfg *config.Config
	db  *pgxpool.Pool
}

func NewPostgresProductRepository(log logger.Logger, cfg *config.Config, db *pgxpool.Pool) *postgresProductRepository {
	return &postgresProductRepository{log: log, cfg: cfg, db: db}
}

func (p *postgresProductRepository) CreateProduct(ctx context.Context, product *models.Product) (*models.Product, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "productRepository.CreateProduct")
	defer span.Finish()

	const createProductQuery = `INSERT INTO products (product_id, name, description, price, created_at, updated_at) 
	VALUES ($1, $2, $3, $4, now(), now()) RETURNING product_id, name, description, price, created_at, updated_at`

	var created models.Product
	if err := p.db.QueryRow(ctx, createProductQuery, &product.ProductID, &product.Name, &product.Description, &product.Price).Scan(
		&created.ProductID,
		&created.Name,
		&created.Description,
		&created.Price,
		&created.CreatedAt,
		&created.UpdatedAt,
	); err != nil {
		return nil, errors.Wrap(err, "error in the insert product into the database")
	}

	return &created, nil
}

func (p *postgresProductRepository) UpdateProduct(ctx context.Context, product *models.Product) (*models.Product, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "productRepository.UpdateProduct")
	defer span.Finish()

	const updateProductQuery = `UPDATE products p SET
	name=COALESCE(NULLIF($1, ''), name),
	description=COALESCE(NULLIF($2, ''), description),
	price=COALESCE(NULLIF($3, 0), price),
	updated_at = now()
	WHERE product_id=$4
	RETURNING product_id, name, description, price, created_at, updated_at`

	var prod models.Product
	if err := p.db.QueryRow(
		ctx,
		updateProductQuery,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.ProductID,
	).Scan(&prod.ProductID, &prod.Name, &prod.Description, &prod.Price, &prod.CreatedAt, &prod.UpdatedAt); err != nil {
		return nil, errors.Wrap(err, "Scan")
	}

	return &prod, nil
}

func (p *postgresProductRepository) GetProductById(ctx context.Context, uuid uuid.UUID) (*models.Product, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "productRepository.GetProductById")
	defer span.Finish()

	const getProductByIdQuery = `SELECT p.product_id, p.name, p.description, p.price, p.created_at, p.updated_at 
	FROM products p WHERE p.product_id = $1`

	var product models.Product
	if err := p.db.QueryRow(ctx, getProductByIdQuery, uuid).Scan(
		&product.ProductID,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.CreatedAt,
		&product.UpdatedAt,
	); err != nil {
		return nil, errors.Wrap(err, "Scan")
	}

	return &product, nil
}

func (p *postgresProductRepository) DeleteProductByID(ctx context.Context, uuid uuid.UUID) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "productRepository.DeleteProductByID")
	defer span.Finish()

	const deleteProductByIdQuery = `DELETE FROM products WHERE product_id = $1`

	_, err := p.db.Exec(ctx, deleteProductByIdQuery, uuid)
	if err != nil {
		return errors.Wrap(err, "Exec")
	}

	return nil
}
