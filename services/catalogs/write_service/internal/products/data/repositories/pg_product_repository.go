package repositories

import (
	"context"
	"fmt"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/tracing"
	"github.com/opentracing/opentracing-go/log"

	"emperror.dev/errors"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/gormPostgres"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/utils"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/models"
	"github.com/opentracing/opentracing-go"
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
		return nil, tracing.TraceWithErr(span, errors.WrapIf(err, "[postgresProductRepository_GetAllProducts.Paginate] error in the paginate"))
	}

	p.log.Infow("[postgresProductRepository.GetAllProducts] products loaded", logger.Fields{"ProductsResult": result})
	span.LogFields(log.Object("ProductsResult", result))

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
		return nil, tracing.TraceWithErr(span, errors.WrapIf(err, "[postgresProductRepository_SearchProducts.Paginate] error in the paginate"))
	}

	p.log.Infow(fmt.Sprintf("[postgresProductRepository.SearchProducts] products loaded for search term '%s'", searchText), logger.Fields{"ProductsResult": result})
	span.LogFields(log.Object("ProductsResult", result))

	return result, nil
}

func (p *postgresProductRepository) GetProductById(ctx context.Context, uuid uuid.UUID) (*models.Product, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "postgresProductRepository.GetProductById")
	span.LogFields(log.String("ProductId", uuid.String()))
	defer span.Finish()

	var product models.Product
	if err := p.gorm.First(&product, uuid).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, tracing.TraceWithErr(span, errors.WrapIf(err, fmt.Sprintf("[postgresProductRepository_GetProductById.First] can't find the product with id %s into the database.", uuid)))
	}
	span.LogFields(log.Object("Product", product))
	p.log.Infow(fmt.Sprintf("[postgresProductRepository.GetProductById] product with id %s laoded", uuid.String()), logger.Fields{"Product": product, "ProductId": uuid})

	return &product, nil
}

func (p *postgresProductRepository) CreateProduct(ctx context.Context, product *models.Product) (*models.Product, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "postgresProductRepository.CreateProduct")
	defer span.Finish()

	if err := p.gorm.Create(&product).Error; err != nil {
		return nil, tracing.TraceWithErr(span, errors.WrapIf(err, "[postgresProductRepository_CreateProduct.Create] error in the inserting product into the database."))
	}
	span.LogFields(log.Object("Product", product))
	p.log.Infow(fmt.Sprintf("[postgresProductRepository.CreateProduct] product with id '%s' created", product.ProductId), logger.Fields{"Product": product, "ProductId": product.ProductId})

	return product, nil
}

func (p *postgresProductRepository) UpdateProduct(ctx context.Context, updateProduct *models.Product) (*models.Product, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "postgresProductRepository.UpdateProduct")
	defer span.Finish()

	if err := p.gorm.Save(updateProduct).Error; err != nil {
		return nil, tracing.TraceWithErr(span, errors.WrapIf(err, fmt.Sprintf("[postgresProductRepository_UpdateProduct.Save] error in updating product with id %s into the database.", updateProduct.ProductId)))
	}
	span.LogFields(log.Object("Product", updateProduct))
	p.log.Infow(fmt.Sprintf("[postgresProductRepository.UpdateProduct] product with id '%s' updated", updateProduct.ProductId), logger.Fields{"Product": updateProduct, "ProductId": updateProduct.ProductId})

	return updateProduct, nil
}

func (p *postgresProductRepository) DeleteProductByID(ctx context.Context, uuid uuid.UUID) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "postgresProductRepository.DeleteProductByID")
	span.LogFields(log.String("AggregateID", uuid.String()))
	defer span.Finish()

	var product models.Product

	if err := p.gorm.First(&product, uuid).Error; err != nil {
		return tracing.TraceWithErr(span, errors.WrapIf(err, fmt.Sprintf("[postgresProductRepository_DeleteProductByID.First] can't find the product with id %s into the database.", uuid)))
	}

	if err := p.gorm.Delete(&product).Error; err != nil {
		return tracing.TraceWithErr(span, errors.WrapIf(err, fmt.Sprintf("[postgresProductRepository_DeleteProductByID.Delete] error in the deleting product with id %s into the database.", uuid)))
	}
	p.log.Infow(fmt.Sprintf("[postgresProductRepository.DeleteProductByID] product with id %s deleted", uuid), logger.Fields{"Product": uuid})

	return nil
}
