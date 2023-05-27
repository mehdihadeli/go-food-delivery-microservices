package repositories

import (
	"context"
	"fmt"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/gorm_postgres/repository"
	"gorm.io/gorm"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/core/data"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/otel/tracing"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/otel/tracing/attribute"
	attribute2 "go.opentelemetry.io/otel/attribute"

	"emperror.dev/errors"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/utils"
	uuid "github.com/satori/go.uuid"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/write_service/internal/products/models"
)

type postgresProductRepository struct {
	log                   logger.Logger
	gormGenericRepository data.GenericRepository[*models.Product]
}

func NewPostgresProductRepository(log logger.Logger, db *gorm.DB) *postgresProductRepository {

	gormRepository := repository.NewGenericGormRepository[*models.Product](db)
	return &postgresProductRepository{log: log, gormGenericRepository: gormRepository}
}

func (p *postgresProductRepository) GetAllProducts(ctx context.Context, listQuery *utils.ListQuery) (*utils.ListResult[*models.Product], error) {
	ctx, span := tracing.Tracer.Start(ctx, "postgresProductRepository.GetAllProducts")
	defer span.End()

	result, err := p.gormGenericRepository.GetAll(ctx, listQuery)
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

	result, err := p.gormGenericRepository.Search(ctx, searchText, listQuery)
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

	product, err := p.gormGenericRepository.GetById(ctx, uuid)
	if err != nil {
		return nil, tracing.TraceErrFromSpan(span, errors.WrapIf(err, fmt.Sprintf("[postgresProductRepository_GetProductById.First] can't find the product with id %s into the database.", uuid)))
	}

	span.SetAttributes(attribute.Object("Product", product))
	p.log.Infow(fmt.Sprintf("[postgresProductRepository.GetProductById] product with id %s laoded", uuid.String()), logger.Fields{"Product": product, "ProductId": uuid})

	return product, nil
}

func (p *postgresProductRepository) CreateProduct(ctx context.Context, product *models.Product) (*models.Product, error) {
	ctx, span := tracing.Tracer.Start(ctx, "postgresProductRepository.CreateProduct")
	defer span.End()

	err := p.gormGenericRepository.Add(ctx, product)
	if err != nil {
		return nil, tracing.TraceErrFromSpan(span, errors.WrapIf(err, "[postgresProductRepository_CreateProduct.Create] error in the inserting product into the database."))
	}

	span.SetAttributes(attribute.Object("Product", product))
	p.log.Infow(fmt.Sprintf("[postgresProductRepository.CreateProduct] product with id '%s' created", product.ProductId), logger.Fields{"Product": product, "ProductId": product.ProductId})

	return product, nil
}

func (p *postgresProductRepository) UpdateProduct(ctx context.Context, updateProduct *models.Product) (*models.Product, error) {
	ctx, span := tracing.Tracer.Start(ctx, "postgresProductRepository.UpdateProduct")
	defer span.End()

	err := p.gormGenericRepository.Update(ctx, updateProduct)
	if err != nil {
		return nil, tracing.TraceErrFromSpan(span, errors.WrapIf(err, fmt.Sprintf("[postgresProductRepository_UpdateProduct.Save] error in updating product with id %s into the database.", updateProduct.ProductId)))
	}

	span.SetAttributes(attribute.Object("Product", updateProduct))
	p.log.Infow(fmt.Sprintf("[postgresProductRepository.UpdateProduct] product with id '%s' updated", updateProduct.ProductId), logger.Fields{"Product": updateProduct, "ProductId": updateProduct.ProductId})

	return updateProduct, nil
}

func (p *postgresProductRepository) DeleteProductByID(ctx context.Context, uuid uuid.UUID) error {
	ctx, span := tracing.Tracer.Start(ctx, "postgresProductRepository.UpdateProduct")
	span.SetAttributes(attribute2.String("ProductId", uuid.String()))
	defer span.End()

	err := p.gormGenericRepository.Delete(ctx, uuid)
	if err != nil {
		return tracing.TraceErrFromSpan(span, errors.WrapIf(err, fmt.Sprintf(
			"[postgresProductRepository_DeleteProductByID.Delete] error in the deleting product with id %s into the database.", uuid)))
	}

	p.log.Infow(fmt.Sprintf("[postgresProductRepository.DeleteProductByID] product with id %s deleted", uuid), logger.Fields{"Product": uuid})

	return nil
}
