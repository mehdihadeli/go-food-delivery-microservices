package repositories

import (
	"context"
	"fmt"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/tracing"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/gorm_postgres"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/utils"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/models"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type postgresProductRepository struct {
	log  logger.Logger
	cfg  *config.Config
	db   *pgxpool.Pool
	gorm *gorm.DB
}

func NewPostgresProductRepository(log logger.Logger, cfg *config.Config, gorm *gorm.DB) *postgresProductRepository {
	return &postgresProductRepository{log: log, cfg: cfg, gorm: gorm}
}

func (p *postgresProductRepository) GetAllProducts(ctx context.Context, listQuery *utils.ListQuery) (*utils.ListResult[models.Product], error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "postgresProductRepository.GetAllProducts")
	defer span.Finish()

	result, err := gorm_postgres.Paginate[models.Product](ctx, listQuery, p.gorm)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	return result, nil
}

func (p *postgresProductRepository) SearchProducts(ctx context.Context, searchText string, listQuery *utils.ListQuery) (*utils.ListResult[models.Product], error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "postgresProductRepository.SearchProducts")
	defer span.Finish()

	whereQuery := fmt.Sprintf("%s IN (?)", "Name")
	query := p.gorm.Where(whereQuery, searchText)

	result, err := gorm_postgres.Paginate[models.Product](ctx, listQuery, query)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	return result, nil
}

func (p *postgresProductRepository) GetProductById(ctx context.Context, uuid uuid.UUID) (*models.Product, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "postgresProductRepository.GetProductById")
	defer span.Finish()

	var product models.Product

	if err := p.gorm.First(&product, uuid).Error; err != nil {
		tracing.TraceErr(span, err)
		return nil, errors.Wrap(err, fmt.Sprintf("can't find the product with id %s into the database.", uuid))
	}

	return &product, nil
}

func (p *postgresProductRepository) CreateProduct(ctx context.Context, product *models.Product) (*models.Product, error) {

	span, ctx := opentracing.StartSpanFromContext(ctx, "postgresProductRepository.CreateProduct")
	defer span.Finish()

	if err := p.gorm.Create(&product).Error; err != nil {
		tracing.TraceErr(span, err)
		return nil, errors.Wrap(err, "error in the inserting product into the database.")
	}

	return product, nil
}

func (p *postgresProductRepository) UpdateProduct(ctx context.Context, updateProduct *models.Product) (*models.Product, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "postgresProductRepository.UpdateProduct")
	defer span.Finish()

	if err := p.gorm.Save(updateProduct).Error; err != nil {
		tracing.TraceErr(span, err)
		return nil, errors.Wrap(err, fmt.Sprintf("error in updating product with id %s into the database.", updateProduct.ProductID))
	}

	return updateProduct, nil
}

func (p *postgresProductRepository) DeleteProductByID(ctx context.Context, uuid uuid.UUID) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "postgresProductRepository.DeleteProductByID")
	defer span.Finish()

	var product models.Product

	if err := p.gorm.First(&product, uuid).Error; err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, fmt.Sprintf("can't find the product with id %s into the database.", uuid))
	}

	if err := p.gorm.Delete(&product).Error; err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "error in the deleting product into the database.")
	}

	return nil
}
