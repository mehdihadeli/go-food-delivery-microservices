package repositories

//https://github.com/Kamva/mgm
//https://github.com/mongodb/mongo-go-driver
//https://blog.logrocket.com/how-to-use-mongodb-with-go/

import (
	"context"
	"emperror.dev/errors"
	"fmt"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/mongodb"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/otel/tracing"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/otel/tracing/attribute"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/utils"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/contracts"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/models"
	uuid "github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	attribute2 "go.opentelemetry.io/otel/attribute"
)

type mongoProductRepository struct {
	log         logger.Logger
	cfg         *config.Config
	mongoClient *mongo.Client
}

func NewMongoProductRepository(log logger.Logger, cfg *config.Config, mongoClient *mongo.Client) contracts.ProductRepository {
	return &mongoProductRepository{log: log, cfg: cfg, mongoClient: mongoClient}
}

func (p *mongoProductRepository) GetAllProducts(ctx context.Context, listQuery *utils.ListQuery) (*utils.ListResult[*models.Product], error) {
	ctx, span := tracing.Tracer.Start(ctx, "mongoProductRepository.GetAllProducts")
	defer span.End()

	collection := p.mongoClient.Database(p.cfg.Mongo.Db).Collection(p.cfg.MongoCollections.Products)

	result, err := mongodb.Paginate[*models.Product](ctx, listQuery, collection, nil)
	if err != nil {
		return nil, tracing.TraceErrFromSpan(span, errors.WrapIf(err, "[mongoProductRepository_GetAllProducts.Paginate] error in the paginate"))
	}

	p.log.Infow("[mongoProductRepository.GetAllProducts] products loaded", logger.Fields{"ProductsResult": result})

	span.SetAttributes(attribute.Object("ProductsResult", result))

	return result, nil
}

func (p *mongoProductRepository) SearchProducts(ctx context.Context, searchText string, listQuery *utils.ListQuery) (*utils.ListResult[*models.Product], error) {
	ctx, span := tracing.Tracer.Start(ctx, "mongoProductRepository.SearchProducts")
	span.SetAttributes(attribute2.String("SearchText", searchText))
	defer span.End()

	collection := p.mongoClient.Database(p.cfg.Mongo.Db).Collection(p.cfg.MongoCollections.Products)

	filter := bson.D{
		{Key: "$or", Value: bson.A{
			bson.D{{Key: "name", Value: primitive.Regex{Pattern: searchText, Options: "gi"}}},
			bson.D{{Key: "description", Value: primitive.Regex{Pattern: searchText, Options: "gi"}}},
		}},
	}

	result, err := mongodb.Paginate[*models.Product](ctx, listQuery, collection, filter)
	if err != nil {
		return nil, tracing.TraceErrFromSpan(span, errors.WrapIf(err, "[mongoProductRepository_SearchProducts.Paginate] error in the paginate"))
	}

	p.log.Infow(fmt.Sprintf("[mongoProductRepository.SearchProducts] products loaded for search term '%s'", searchText), logger.Fields{"ProductsResult": result})

	span.SetAttributes(attribute.Object("ProductsResult", result))

	return result, nil
}

func (p *mongoProductRepository) GetProductById(ctx context.Context, uuid uuid.UUID) (*models.Product, error) {
	ctx, span := tracing.Tracer.Start(ctx, "mongoProductRepository.GetProductById")
	span.SetAttributes(attribute2.String("Id", uuid.String()))
	defer span.End()

	collection := p.mongoClient.Database(p.cfg.Mongo.Db).Collection(p.cfg.MongoCollections.Products)

	var product models.Product
	if err := collection.FindOne(ctx, bson.M{"_id": uuid.String()}).Decode(&product); err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, tracing.TraceErrFromSpan(span, errors.WrapIf(err, fmt.Sprintf("[mongoProductRepository_GetProductById.FindOne] can't find the product with id %s into the database.", uuid)))
	}

	span.SetAttributes(attribute.Object("Product", product))

	p.log.Infow(fmt.Sprintf("[mongoProductRepository.GetProductById] product with id %s laoded", uuid.String()), logger.Fields{"Product": product, "Id": uuid})

	return &product, nil
}

func (p *mongoProductRepository) GetProductByProductId(ctx context.Context, uuid uuid.UUID) (*models.Product, error) {
	ctx, span := tracing.Tracer.Start(ctx, "mongoProductRepository.GetProductByProductId")
	span.SetAttributes(attribute2.String("ProductId", uuid.String()))
	defer span.End()

	collection := p.mongoClient.Database(p.cfg.Mongo.Db).Collection(p.cfg.MongoCollections.Products)

	var product models.Product
	if err := collection.FindOne(ctx, bson.M{"productId": uuid.String()}).Decode(&product); err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, tracing.TraceErrFromSpan(span, errors.WrapIf(err, fmt.Sprintf("[mongoProductRepository_GetProductById.FindOne] can't find the product with productId %s into the database.", uuid)))
	}

	span.SetAttributes(attribute.Object("Product", product))

	p.log.Infow(fmt.Sprintf("[mongoProductRepository.GetProductById] product with productId %s laoded", uuid.String()), logger.Fields{"Product": product, "Id": uuid})

	return &product, nil
}

func (p *mongoProductRepository) CreateProduct(ctx context.Context, product *models.Product) (*models.Product, error) {
	ctx, span := tracing.Tracer.Start(ctx, "mongoProductRepository.CreateProduct")
	defer span.End()

	collection := p.mongoClient.Database(p.cfg.Mongo.Db).Collection(p.cfg.MongoCollections.Products)
	_, err := collection.InsertOne(ctx, product, &options.InsertOneOptions{})
	if err != nil {
		return nil, tracing.TraceErrFromSpan(span, errors.WrapIf(err, "[mongoProductRepository_CreateProduct.InsertOne] error in the inserting product into the database."))
	}

	span.SetAttributes(attribute.Object("Product", product))

	p.log.Infow(fmt.Sprintf("[mongoProductRepository.CreateProduct] product with id '%s' created", product.ProductId), logger.Fields{"Product": product, "Id": product.ProductId})

	return product, nil
}

func (p *mongoProductRepository) UpdateProduct(ctx context.Context, updateProduct *models.Product) (*models.Product, error) {
	ctx, span := tracing.Tracer.Start(ctx, "mongoProductRepository.UpdateProduct")
	defer span.End()

	collection := p.mongoClient.Database(p.cfg.Mongo.Db).Collection(p.cfg.MongoCollections.Products)

	ops := options.FindOneAndUpdate()
	ops.SetReturnDocument(options.After)
	ops.SetUpsert(true)

	var updated models.Product
	//https://www.mongodb.com/docs/manual/reference/method/db.collection.findOneAndUpdate/
	if err := collection.FindOneAndUpdate(ctx, bson.M{"_id": updateProduct.Id}, bson.M{"$set": updateProduct}, ops).Decode(&updated); err != nil {
		return nil, tracing.TraceErrFromSpan(span, errors.WrapIf(err, fmt.Sprintf("[mongoProductRepository_UpdateProduct.FindOneAndUpdate] error in updating product with id %s into the database.", updateProduct.ProductId)))
	}

	span.SetAttributes(attribute.Object("Product", updateProduct))
	p.log.Infow(fmt.Sprintf("[mongoProductRepository.UpdateProduct] product with id '%s' updated", updateProduct.ProductId), logger.Fields{"Product": updateProduct, "Id": updateProduct.ProductId})

	return &updated, nil
}

func (p *mongoProductRepository) DeleteProductByID(ctx context.Context, uuid uuid.UUID) error {
	ctx, span := tracing.Tracer.Start(ctx, "mongoProductRepository.DeleteProductByID")
	span.SetAttributes(attribute2.String("Id", uuid.String()))
	defer span.End()

	collection := p.mongoClient.Database(p.cfg.Mongo.Db).Collection(p.cfg.MongoCollections.Products)

	if err := collection.FindOneAndDelete(ctx, bson.M{"_id": uuid.String()}).Err(); err != nil {
		return tracing.TraceErrFromSpan(span, errors.WrapIf(err, fmt.Sprintf(
			"[mongoProductRepository_DeleteProductByID.FindOneAndDelete] error in deleting product with id %d from the database.", uuid)))
	}

	p.log.Infow(fmt.Sprintf("[mongoProductRepository.DeleteProductByID] product with id %s deleted", uuid), logger.Fields{"Product": uuid})

	return nil
}
