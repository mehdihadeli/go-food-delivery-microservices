package repositories

import (
	"context"
	"fmt"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/otel/tracing"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/otel/tracing/attribute"
	attribute2 "go.opentelemetry.io/otel/attribute"

	"emperror.dev/errors"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/gormPostgres"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/utils"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/models"
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
	ctx, span := tracing.Tracer.Start(ctx, "postgresProductRepository.GetAllProducts")
	defer span.End()

	result, err := gormPostgres.Paginate[*models.Product](ctx, listQuery, p.gorm)
	if err != nil {
		return nil, tracing.TraceErrFromContext(ctx, errors.WrapIf(err, "[postgresProductRepository_GetAllProducts.Paginate] error in the paginate"))
	}

	p.log.Infow("[postgresProductRepository.GetAllProducts] products loaded", logger.Fields{"ProductsResult": result})
	span.SetAttributes(attribute.Object("ProductsResult", result))

	return result, nil
}

func (p *postgresProductRepository) SearchProducts(ctx context.Context, searchText string, listQuery *utils.ListQuery) (*utils.ListResult[*models.Product], error) {
	ctx, span := tracing.Tracer.Start(ctx, "postgresProductRepository.SearchProducts")
	span.SetAttributes(attribute2.String("SearchText", searchText))
	defer span.End()

	whereQuery := fmt.Sprintf("%s IN (?)", "Name")
	query := p.gorm.Where(whereQuery, searchText)

	result, err := gormPostgres.Paginate[*models.Product](ctx, listQuery, query)
	if err != nil {
		return nil, tracing.TraceErrFromContext(ctx, errors.WrapIf(err, "[postgresProductRepository_SearchProducts.Paginate] error in the paginate"))
	}

	p.log.Infow(fmt.Sprintf("[postgresProductRepository.SearchProducts] products loaded for search term '%s'", searchText), logger.Fields{"ProductsResult": result})
	span.SetAttributes(attribute.Object("ProductsResult", result))

	return result, nil
}

func (p *postgresProductRepository) GetProductById(ctx context.Context, uuid uuid.UUID) (*models.Product, error) {
	ctx, span := tracing.Tracer.Start(ctx, "postgresProductRepository.GetProductById")
	span.SetAttributes(attribute2.String("ProductId", uuid.String()))
	defer span.End()

	var product models.Product
	if err := p.gorm.First(&product, uuid).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, tracing.TraceErrFromContext(ctx, errors.WrapIf(err, fmt.Sprintf("[postgresProductRepository_GetProductById.First] can't find the product with id %s into the database.", uuid)))
	}
	span.SetAttributes(attribute.Object("Product", product))
	p.log.Infow(fmt.Sprintf("[postgresProductRepository.GetProductById] product with id %s laoded", uuid.String()), logger.Fields{"Product": product, "ProductId": uuid})

	return &product, nil
}

func (p *postgresProductRepository) CreateProduct(ctx context.Context, product *models.Product) (*models.Product, error) {
	ctx, span := tracing.Tracer.Start(ctx, "postgresProductRepository.CreateProduct")
	defer span.End()

	if err := p.gorm.Create(&product).Error; err != nil {
		return nil, tracing.TraceErrFromContext(ctx, errors.WrapIf(err, "[postgresProductRepository_CreateProduct.Create] error in the inserting product into the database."))
	}
	span.SetAttributes(attribute.Object("Product", product))
	p.log.Infow(fmt.Sprintf("[postgresProductRepository.CreateProduct] product with id '%s' created", product.ProductId), logger.Fields{"Product": product, "ProductId": product.ProductId})

	return product, nil
}

func (p *postgresProductRepository) UpdateProduct(ctx context.Context, updateProduct *models.Product) (*models.Product, error) {
	ctx, span := tracing.Tracer.Start(ctx, "postgresProductRepository.UpdateProduct")
	defer span.End()

	if err := p.gorm.Save(updateProduct).Error; err != nil {
		return nil, tracing.TraceErrFromContext(ctx, errors.WrapIf(err, fmt.Sprintf("[postgresProductRepository_UpdateProduct.Save] error in updating product with id %s into the database.", updateProduct.ProductId)))
	}
	span.SetAttributes(attribute.Object("Product", updateProduct))
	p.log.Infow(fmt.Sprintf("[postgresProductRepository.UpdateProduct] product with id '%s' updated", updateProduct.ProductId), logger.Fields{"Product": updateProduct, "ProductId": updateProduct.ProductId})

	return updateProduct, nil
}

func (p *postgresProductRepository) DeleteProductByID(ctx context.Context, uuid uuid.UUID) error {
	ctx, span := tracing.Tracer.Start(ctx, "postgresProductRepository.UpdateProduct")
	span.SetAttributes(attribute2.String("ProductId", uuid.String()))
	defer span.End()

	var product models.Product

	if err := p.gorm.First(&product, uuid).Error; err != nil {
		return tracing.TraceErrFromContext(ctx, errors.WrapIf(err, fmt.Sprintf("[postgresProductRepository_DeleteProductByID.First] can't find the product with id %s into the database.", uuid)))
	}

	if err := p.gorm.Delete(&product).Error; err != nil {
		return tracing.TraceErrFromContext(ctx, errors.WrapIf(err, fmt.Sprintf("[postgresProductRepository_DeleteProductByID.Delete] error in the deleting product with id %s into the database.", uuid)))
	}
	p.log.Infow(fmt.Sprintf("[postgresProductRepository.DeleteProductByID] product with id %s deleted", uuid), logger.Fields{"Product": uuid})

	return nil
}
