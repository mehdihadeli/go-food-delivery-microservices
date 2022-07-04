package repositories

import (
	"context"
	"fmt"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/gorm_postgres"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/tracing"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/utils"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/models"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
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

func (p *mongoProductRepository) GetAllProducts(ctx context.Context, listQuery *utils.ListQuery) (*utils.ListResult[models.Product], error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "mongoProductRepository.GetAllProducts")
	defer span.Finish()

	result, err := gorm_postgres.Paginate[models.Product](listQuery, p.gorm)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (p *mongoProductRepository) SearchProducts(ctx context.Context, searchText string, listQuery *utils.ListQuery) (*utils.ListResult[models.Product], error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "mongoProductRepository.SearchProducts")
	defer span.Finish()

	whereQuery := fmt.Sprintf("%s IN (?)", "Name")
	query := p.gorm.Where(whereQuery, searchText)

	result, err := gorm_postgres.Paginate[models.Product](listQuery, query)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (p *mongoProductRepository) GetProductById(ctx context.Context, uuid uuid.UUID) (*models.Product, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "mongoProductRepository.GetProductById")
	defer span.Finish()

	var product models.Product

	if result := p.gorm.First(&product, uuid); result.Error != nil {
		return nil, errors.Wrap(result.Error, fmt.Sprintf("can't find the product with id %s into the database.", uuid))
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

	if result := p.gorm.Save(updateProduct); result.Error != nil {
		return nil, errors.Wrap(result.Error, fmt.Sprintf("error in updating product with id %s into the database.", updateProduct.ProductID))
	}

	return updateProduct, nil
}

func (p *mongoProductRepository) DeleteProductByID(ctx context.Context, uuid uuid.UUID) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "mongoProductRepository.DeleteProductByID")
	defer span.Finish()

	var product models.Product

	if result := p.gorm.First(&product, uuid); result.Error != nil {
		return errors.Wrap(result.Error, fmt.Sprintf("can't find the product with id %s into the database.", uuid))
	}

	if result := p.gorm.Delete(&product); result.Error != nil {
		return errors.Wrap(result.Error, "error in the deleting product into the database.")
	}

	return nil
}
