package repositories

//https://github.com/Kamva/mgm
//https://github.com/mongodb/mongo-go-driver
//https://blog.logrocket.com/how-to-use-mongodb-with-go/

import (
	"context"
	"fmt"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/mongodb"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/tracing"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/utils"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/models"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoProductRepository struct {
	log         logger.Logger
	cfg         *config.Config
	mongoClient *mongo.Client
}

func NewMongoProductRepository(log logger.Logger, cfg *config.Config, mongoClient *mongo.Client) *mongoProductRepository {
	return &mongoProductRepository{log: log, cfg: cfg, mongoClient: mongoClient}
}

func (p *mongoProductRepository) GetAllProducts(ctx context.Context, listQuery *utils.ListQuery) (*utils.ListResult[*models.Product], error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "mongoProductRepository.GetAllProducts")
	defer span.Finish()

	collection := p.mongoClient.Database(p.cfg.Mongo.Db).Collection(p.cfg.MongoCollections.Products)

	result, err := mongodb.Paginate[*models.Product](ctx, listQuery, collection, nil)
	if err != nil {
		return nil, tracing.TraceWithErr(span, errors.Wrap(err, "[mongoProductRepository_GetAllProducts.Paginate] error in the paginate"))
	}

	p.log.Info("[mongoProductRepository.GetAllProducts] result: %+v", result)

	return result, nil
}

func (p *mongoProductRepository) SearchProducts(ctx context.Context, searchText string, listQuery *utils.ListQuery) (*utils.ListResult[*models.Product], error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "mongoProductRepository.SearchProducts")
	span.LogFields(log.String("SearchText", searchText))
	defer span.Finish()

	collection := p.mongoClient.Database(p.cfg.Mongo.Db).Collection(p.cfg.MongoCollections.Products)

	filter := bson.D{
		{Key: "$or", Value: bson.A{
			bson.D{{Key: "name", Value: primitive.Regex{Pattern: searchText, Options: "gi"}}},
			bson.D{{Key: "description", Value: primitive.Regex{Pattern: searchText, Options: "gi"}}},
		}},
	}

	result, err := mongodb.Paginate[*models.Product](ctx, listQuery, collection, filter)
	if err != nil {
		return nil, tracing.TraceWithErr(span, errors.Wrap(err, "[mongoProductRepository_SearchProducts.Paginate] error in the paginate"))
	}

	p.log.Info("[mongoProductRepository.SearchProducts] result: %+v", result)
	return result, nil
}

func (p *mongoProductRepository) GetProductById(ctx context.Context, uuid uuid.UUID) (*models.Product, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "mongoProductRepository.GetProductById")
	span.LogFields(log.String("AggregateID", uuid.String()))
	defer span.Finish()

	collection := p.mongoClient.Database(p.cfg.Mongo.Db).Collection(p.cfg.MongoCollections.Products)

	var product models.Product
	if err := collection.FindOne(ctx, bson.M{"_id": uuid.String()}).Decode(&product); err != nil {
		return nil, tracing.TraceWithErr(span, errors.Wrap(err, fmt.Sprintf("[mongoProductRepository_GetProductById.FindOne] can't find the product with id %s into the database.", uuid)))
	}

	p.log.Infow(fmt.Sprintf("[mongoProductRepository.GetProductById] result: %+v", product), logger.Fields{"AggregateID": uuid})
	return &product, nil
}

func (p *mongoProductRepository) CreateProduct(ctx context.Context, product *models.Product) (*models.Product, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "mongoProductRepository.CreateProduct")
	span.LogFields(log.Object("Aggregate", product))
	defer span.Finish()

	collection := p.mongoClient.Database(p.cfg.Mongo.Db).Collection(p.cfg.MongoCollections.Products)
	_, err := collection.InsertOne(ctx, product, &options.InsertOneOptions{})
	if err != nil {
		return nil, tracing.TraceWithErr(span, errors.Wrap(err, "[mongoProductRepository_CreateProduct.InsertOne] error in the inserting product into the database."))
	}

	p.log.Infow(fmt.Sprintf("[mongoProductRepository.CreateProduct] result AggregateID: %s", product.ProductID), logger.Fields{"AggregateID": product.ProductID})
	return product, nil
}

func (p *mongoProductRepository) UpdateProduct(ctx context.Context, updateProduct *models.Product) (*models.Product, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "mongoProductRepository.UpdateProduct")
	span.LogFields(log.Object("Aggregate", updateProduct))
	defer span.Finish()

	collection := p.mongoClient.Database(p.cfg.Mongo.Db).Collection(p.cfg.MongoCollections.Products)

	ops := options.FindOneAndUpdate()
	ops.SetReturnDocument(options.After)
	ops.SetUpsert(true)

	var updated models.Product
	if err := collection.FindOneAndUpdate(ctx, bson.M{"_id": updateProduct.ProductID}, bson.M{"$set": updateProduct}, ops).Decode(&updated); err != nil {
		return nil, tracing.TraceWithErr(span, errors.Wrap(err, fmt.Sprintf("[mongoProductRepository_UpdateProduct.FindOneAndUpdate] error in updating product with id %s into the database.", updateProduct.ProductID)))
	}

	p.log.Infow(fmt.Sprintf("[mongoProductRepository.UpdateProduct] result AggregateID: %s", updateProduct.ProductID), logger.Fields{"AggregateID": updateProduct.ProductID})
	return &updated, nil
}

func (p *mongoProductRepository) DeleteProductByID(ctx context.Context, uuid uuid.UUID) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "mongoProductRepository.DeleteProductByID")
	span.LogFields(log.String("AggregateID", uuid.String()))
	defer span.Finish()

	collection := p.mongoClient.Database(p.cfg.Mongo.Db).Collection(p.cfg.MongoCollections.Products)

	if err := collection.FindOneAndDelete(ctx, bson.M{"_id": uuid.String()}).Err(); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, fmt.Sprintf(
			"[mongoProductRepository_DeleteProductByID.FindOneAndDelete] error in deleting product with id %d from the database.", uuid)))
	}

	p.log.Infow(fmt.Sprintf("[mongoProductRepository.DeleteProductByID] result AggregateID: %s", uuid), logger.Fields{"AggregateID": uuid})
	return nil
}
