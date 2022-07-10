package repositories

//https://github.com/Kamva/mgm
//https://github.com/mongodb/mongo-go-driver
//https://blog.logrocket.com/how-to-use-mongodb-with-go/

import (
	"context"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/mongodb"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/tracing"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/utils"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/models"
	"github.com/opentracing/opentracing-go"
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
		tracing.TraceErr(span, err)
		return nil, err
	}

	return result, nil
}

func (p *mongoProductRepository) SearchProducts(ctx context.Context, searchText string, listQuery *utils.ListQuery) (*utils.ListResult[*models.Product], error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "mongoProductRepository.SearchProducts")
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
		tracing.TraceErr(span, err)
		return nil, err
	}

	return result, nil
}

func (p *mongoProductRepository) GetProductById(ctx context.Context, uuid uuid.UUID) (*models.Product, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "mongoProductRepository.GetProductById")
	defer span.Finish()

	collection := p.mongoClient.Database(p.cfg.Mongo.Db).Collection(p.cfg.MongoCollections.Products)

	var product models.Product
	if err := collection.FindOne(ctx, bson.M{"_id": uuid.String()}).Decode(&product); err != nil {
		tracing.TraceErr(span, err)
		return nil, errors.Wrap(err, "error in the getting product from the database.")
	}

	return &product, nil
}

func (p *mongoProductRepository) CreateProduct(ctx context.Context, product *models.Product) (*models.Product, error) {

	span, ctx := opentracing.StartSpanFromContext(ctx, "mongoProductRepository.CreateProduct")
	defer span.Finish()

	collection := p.mongoClient.Database(p.cfg.Mongo.Db).Collection(p.cfg.MongoCollections.Products)

	_, err := collection.InsertOne(ctx, product, &options.InsertOneOptions{})
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, errors.Wrap(err, "error in the inserting product into the database.")
	}

	return product, nil
}

func (p *mongoProductRepository) UpdateProduct(ctx context.Context, updateProduct *models.Product) (*models.Product, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "mongoProductRepository.UpdateProduct")
	defer span.Finish()

	collection := p.mongoClient.Database(p.cfg.Mongo.Db).Collection(p.cfg.MongoCollections.Products)

	ops := options.FindOneAndUpdate()
	ops.SetReturnDocument(options.After)
	ops.SetUpsert(true)

	var updated models.Product
	if err := collection.FindOneAndUpdate(ctx, bson.M{"_id": updateProduct.ProductID}, bson.M{"$set": updateProduct}, ops).Decode(&updated); err != nil {
		tracing.TraceErr(span, err)
		return nil, errors.Wrap(err, "error in the updating product into the database.")
	}

	return &updated, nil
}

func (p *mongoProductRepository) DeleteProductByID(ctx context.Context, uuid uuid.UUID) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "mongoProductRepository.DeleteProductByID")
	defer span.Finish()

	collection := p.mongoClient.Database(p.cfg.Mongo.Db).Collection(p.cfg.MongoCollections.Products)

	if err := collection.FindOneAndDelete(ctx, bson.M{"_id": uuid.String()}).Err(); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "error in deleting product from the database.")
	}

	return nil
}
