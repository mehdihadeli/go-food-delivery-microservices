package repositories

import (
	"context"
	"fmt"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/tracing"
	"github.com/opentracing/opentracing-go/log"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/gormPostgres"
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

func (p *postgresProductRepository) GetAllProducts(ctx context.Context, listQuery *utils.ListQuery) (*utils.ListResult[*models.Product], error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "postgresProductRepository.GetAllProducts")
	defer span.Finish()

	result, err := gormPostgres.Paginate[*models.Product](ctx, listQuery, p.gorm)
	if err != nil {
		return nil, tracing.TraceWithErr(span, errors.Wrap(err, "[postgresProductRepository_GetAllProducts.Paginate] error in the paginate"))
	}

	p.log.Info("[postgresProductRepository.GetAllProducts] result: %+v", result)
	return result, nil
}

func (p *postgresProductRepository) SearchProducts(ctx context.Context, searchText string, listQuery *utils.ListQuery) (*utils.ListResult[*models.Product], error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "postgresProductRepository.SearchProducts")
	span.LogFields(log.String("SearchText", searchText))
	defer span.Finish()

	whereQuery := fmt.Sprintf("%s IN (?)", "Name")
	query := p.gorm.Where(whereQuery, searchText)

	result, err := gormPostgres.Paginate[*models.Product](ctx, listQuery, query)
	if err != nil {
		return nil, tracing.TraceWithErr(span, errors.Wrap(err, "[postgresProductRepository_SearchProducts.Paginate] error in the paginate"))
	}

	p.log.Info("[postgresProductRepository.SearchProducts] result: %+v", result)
	return result, nil
}

func (p *postgresProductRepository) GetProductById(ctx context.Context, uuid uuid.UUID) (*models.Product, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "postgresProductRepository.GetProductById")
	span.LogFields(log.String("AggregateID", uuid.String()))
	defer span.Finish()

	var product models.Product
	if err := p.gorm.First(&product, uuid).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, tracing.TraceWithErr(span, errors.Wrap(err, fmt.Sprintf("[postgresProductRepository_GetProductById.First] can't find the product with id %s into the database.", uuid)))
	}

	p.log.Infow(fmt.Sprintf("[postgresProductRepository.GetProductById] result: %+v", product), logger.Fields{"AggregateID": uuid})
	return &product, nil
}

func (p *postgresProductRepository) CreateProduct(ctx context.Context, product *models.Product) (*models.Product, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "postgresProductRepository.CreateProduct")
	span.LogFields(log.Object("Aggregate", product))
	defer span.Finish()

	if err := p.gorm.Create(&product).Error; err != nil {
		return nil, tracing.TraceWithErr(span, errors.Wrap(err, "[postgresProductRepository_CreateProduct.Create] error in the inserting product into the database."))
	}

	p.log.Infow(fmt.Sprintf("[postgresProductRepository.CreateProduct] result AggregateID: %s", product.ProductID), logger.Fields{"AggregateID": product.ProductID})
	return product, nil
}

func (p *postgresProductRepository) UpdateProduct(ctx context.Context, updateProduct *models.Product) (*models.Product, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "postgresProductRepository.UpdateProduct")
	span.LogFields(log.Object("Aggregate", updateProduct))
	defer span.Finish()

	if err := p.gorm.Save(updateProduct).Error; err != nil {
		return nil, tracing.TraceWithErr(span, errors.Wrap(err, fmt.Sprintf("[postgresProductRepository_UpdateProduct.Save] error in updating product with id %s into the database.", updateProduct.ProductID)))
	}

	p.log.Infow(fmt.Sprintf("[postgresProductRepository.UpdateProduct] result AggregateID: %s", updateProduct.ProductID), logger.Fields{"AggregateID": updateProduct.ProductID})
	return updateProduct, nil
}

func (p *postgresProductRepository) DeleteProductByID(ctx context.Context, uuid uuid.UUID) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "postgresProductRepository.DeleteProductByID")
	span.LogFields(log.String("AggregateID", uuid.String()))
	defer span.Finish()

	var product models.Product

	if err := p.gorm.First(&product, uuid).Error; err != nil {
		return tracing.TraceWithErr(span, errors.Wrap(err, fmt.Sprintf("[postgresProductRepository_DeleteProductByID.First] can't find the product with id %s into the database.", uuid)))
	}

	if err := p.gorm.Delete(&product).Error; err != nil {
		return tracing.TraceWithErr(span, errors.Wrap(err, fmt.Sprintf("[postgresProductRepository_DeleteProductByID.Delete] error in the deleting product with id %s into the database.", uuid)))
	}

	p.log.Infow(fmt.Sprintf("[postgresProductRepository.DeleteProductByID] result AggregateID: %s", uuid), logger.Fields{"AggregateID": uuid})
	return nil
}
