package repositories

// https://github.com/Kamva/mgm
// https://github.com/mongodb/mongo-go-driver
// https://blog.logrocket.com/how-to-use-mongodb-with-go/
// https://www.mongodb.com/docs/drivers/go/current/quick-reference/
// https://www.mongodb.com/docs/drivers/go/current/fundamentals/bson/
// https://www.mongodb.com/docs

import (
	"context"
	"fmt"

	"emperror.dev/errors"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/core/data"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/mongodb"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/mongodb/repository"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/otel/tracing"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/otel/tracing/attribute"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/utils"
	uuid2 "github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/mongo"
	attribute2 "go.opentelemetry.io/otel/attribute"

	data2 "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogreadservice/internal/products/contracts/data"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogreadservice/internal/products/models"
)

const (
	productCollection = "products"
)

type mongoProductRepository struct {
	log                    logger.Logger
	mongoGenericRepository data.GenericRepository[*models.Product]
	tracer                 tracing.AppTracer
}

func NewMongoProductRepository(
	log logger.Logger,
	db *mongo.Client,
	mongoOptions *mongodb.MongoDbOptions,
	tracer tracing.AppTracer,
) data2.ProductRepository {
	mongoRepo := repository.NewGenericMongoRepository[*models.Product](
		db,
		mongoOptions.Database,
		productCollection,
	)
	return &mongoProductRepository{log: log, mongoGenericRepository: mongoRepo, tracer: tracer}
}

func (p *mongoProductRepository) GetAllProducts(
	ctx context.Context,
	listQuery *utils.ListQuery,
) (*utils.ListResult[*models.Product], error) {
	ctx, span := p.tracer.Start(ctx, "mongoProductRepository.GetAllProducts")
	defer span.End()

	// https://www.mongodb.com/docs/drivers/go/current/fundamentals/crud/read-operations/query-document/
	result, err := p.mongoGenericRepository.GetAll(ctx, listQuery)
	if err != nil {
		return nil, tracing.TraceErrFromSpan(
			span,
			errors.WrapIf(
				err,
				"[mongoProductRepository_GetAllProducts.Paginate] error in the paginate",
			),
		)
	}

	p.log.Infow(
		"[mongoProductRepository.GetAllProducts] products loaded",
		logger.Fields{"ProductsResult": result},
	)

	span.SetAttributes(attribute.Object("ProductsResult", result))

	return result, nil
}

func (p *mongoProductRepository) SearchProducts(
	ctx context.Context,
	searchText string,
	listQuery *utils.ListQuery,
) (*utils.ListResult[*models.Product], error) {
	ctx, span := p.tracer.Start(ctx, "mongoProductRepository.SearchProducts")
	span.SetAttributes(attribute2.String("SearchText", searchText))
	defer span.End()

	result, err := p.mongoGenericRepository.Search(ctx, searchText, listQuery)
	if err != nil {
		return nil, tracing.TraceErrFromSpan(
			span,
			errors.WrapIf(
				err,
				"[mongoProductRepository_SearchProducts.Paginate] error in the paginate",
			),
		)
	}

	p.log.Infow(
		fmt.Sprintf(
			"[mongoProductRepository.SearchProducts] products loaded for search term '%s'",
			searchText,
		),
		logger.Fields{"ProductsResult": result},
	)

	span.SetAttributes(attribute.Object("ProductsResult", result))

	return result, nil
}

func (p *mongoProductRepository) GetProductById(
	ctx context.Context,
	uuid string,
) (*models.Product, error) {
	ctx, span := p.tracer.Start(ctx, "mongoProductRepository.GetProductById")
	span.SetAttributes(attribute2.String("Id", uuid))
	defer span.End()

	id, err := uuid2.FromString(uuid)
	if err != nil {
		return nil, err
	}

	product, err := p.mongoGenericRepository.GetById(ctx, id)
	if err != nil {
		return nil, tracing.TraceErrFromSpan(
			span,
			errors.WrapIf(
				err,
				fmt.Sprintf(
					"[mongoProductRepository_GetProductById.FindOne] can't find the product with id %s into the database.",
					uuid,
				),
			),
		)
	}

	span.SetAttributes(attribute.Object("Product", product))

	p.log.Infow(
		fmt.Sprintf("[mongoProductRepository.GetProductById] product with id %s laoded", uuid),
		logger.Fields{"Product": product, "Id": uuid},
	)

	return product, nil
}

func (p *mongoProductRepository) GetProductByProductId(
	ctx context.Context,
	uuid string,
) (*models.Product, error) {
	productId := uuid
	ctx, span := p.tracer.Start(ctx, "mongoProductRepository.GetProductByProductId")
	span.SetAttributes(attribute2.String("ProductId", productId))
	defer span.End()

	product, err := p.mongoGenericRepository.FirstOrDefault(
		ctx,
		map[string]interface{}{"productId": uuid},
	)
	if err != nil {
		return nil, tracing.TraceErrFromSpan(
			span,
			errors.WrapIf(
				err,
				fmt.Sprintf(
					"[mongoProductRepository_GetProductById.FindOne] can't find the product with productId %s into the database.",
					uuid,
				),
			),
		)
	}

	span.SetAttributes(attribute.Object("Product", product))

	p.log.Infow(
		fmt.Sprintf(
			"[mongoProductRepository.GetProductById] product with productId %s laoded",
			productId,
		),
		logger.Fields{"Product": product, "ProductId": uuid},
	)

	return product, nil
}

func (p *mongoProductRepository) CreateProduct(
	ctx context.Context,
	product *models.Product,
) (*models.Product, error) {
	ctx, span := p.tracer.Start(ctx, "mongoProductRepository.CreateProduct")
	defer span.End()

	err := p.mongoGenericRepository.Add(ctx, product)
	if err != nil {
		return nil, tracing.TraceErrFromSpan(
			span,
			errors.WrapIf(
				err,
				"[mongoProductRepository_CreateProduct.InsertOne] error in the inserting product into the database.",
			),
		)
	}

	span.SetAttributes(attribute.Object("Product", product))

	p.log.Infow(
		fmt.Sprintf(
			"[mongoProductRepository.CreateProduct] product with id '%s' created",
			product.ProductId,
		),
		logger.Fields{"Product": product, "Id": product.ProductId},
	)

	return product, nil
}

func (p *mongoProductRepository) UpdateProduct(
	ctx context.Context,
	updateProduct *models.Product,
) (*models.Product, error) {
	ctx, span := p.tracer.Start(ctx, "mongoProductRepository.UpdateProduct")
	defer span.End()

	err := p.mongoGenericRepository.Update(ctx, updateProduct)
	// https://www.mongodb.com/docs/manual/reference/method/db.collection.findOneAndUpdate/
	if err != nil {
		return nil, tracing.TraceErrFromSpan(
			span,
			errors.WrapIf(
				err,
				fmt.Sprintf(
					"[mongoProductRepository_UpdateProduct.FindOneAndUpdate] error in updating product with id %s into the database.",
					updateProduct.ProductId,
				),
			),
		)
	}

	span.SetAttributes(attribute.Object("Product", updateProduct))
	p.log.Infow(
		fmt.Sprintf(
			"[mongoProductRepository.UpdateProduct] product with id '%s' updated",
			updateProduct.ProductId,
		),
		logger.Fields{"Product": updateProduct, "Id": updateProduct.ProductId},
	)

	return updateProduct, nil
}

func (p *mongoProductRepository) DeleteProductByID(ctx context.Context, uuid string) error {
	ctx, span := p.tracer.Start(ctx, "mongoProductRepository.DeleteProductByID")
	span.SetAttributes(attribute2.String("Id", uuid))
	defer span.End()

	id, err := uuid2.FromString(uuid)
	if err != nil {
		return err
	}

	err = p.mongoGenericRepository.Delete(ctx, id)
	if err != nil {
		return tracing.TraceErrFromSpan(span, errors.WrapIf(err, fmt.Sprintf(
			"[mongoProductRepository_DeleteProductByID.FindOneAndDelete] error in deleting product with id %s from the database.",
			uuid,
		)))
	}

	p.log.Infow(
		fmt.Sprintf("[mongoProductRepository.DeleteProductByID] product with id %s deleted", uuid),
		logger.Fields{"Product": uuid},
	)

	return nil
}
