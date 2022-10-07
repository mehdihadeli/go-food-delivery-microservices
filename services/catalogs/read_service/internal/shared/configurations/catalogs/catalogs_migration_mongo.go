package catalogs

import (
	"context"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/consts"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (c *catalogsServiceConfigurator) migrationCatalogsMongo(ctx context.Context, mongoClient *mongo.Client) {
	err := mongoClient.Database(c.GetCfg().Mongo.Db).CreateCollection(ctx, c.GetCfg().MongoCollections.Products)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			c.GetLog().Warnf("(CreateCollection) err: {%v}", err)
		}
	}

	indexOptions := options.Index().SetSparse(true).SetUnique(true)
	index, err := mongoClient.Database(c.GetCfg().Mongo.Db).Collection(c.GetCfg().MongoCollections.Products).Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: consts.ProductIdIndex, Value: 1}},
		Options: indexOptions,
	})
	if err != nil && mongo.IsDuplicateKeyError(err) {
		c.GetLog().Warnf("(CreateOne) err: {%v}", err)
	}
	c.GetLog().Infof("(CreatedIndex) index: {%s}", index)

	list, err := mongoClient.Database(c.GetCfg().Mongo.Db).Collection(c.GetCfg().MongoCollections.Products).Indexes().List(ctx)
	if err != nil {
		c.GetLog().Warnf("(initMongoDBCollections) [List] err: {%v}", err)
	}

	if list != nil {
		var results []bson.M
		if err := list.All(ctx, &results); err != nil {
			c.GetLog().Warnf("(All) err: {%v}", err)
		}
		c.GetLog().Infof("(indexes) results: {%#v}", results)
	}

	collections, err := mongoClient.Database(c.GetCfg().Mongo.Db).ListCollectionNames(ctx, bson.M{})
	if err != nil {
		c.GetLog().Warnf("(ListCollections) err: {%v}", err)
	}
	c.GetLog().Infof("(Collections) created collections: {%v}", collections)
}
