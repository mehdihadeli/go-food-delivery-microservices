package catalogs

import (
	"context"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/utils"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/consts"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/shared/constants"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (c *catalogsServiceConfigurator) migrationCatalogsMongo(ctx context.Context, mongoClient *mongo.Client) {
	err := mongoClient.Database(c.cfg.Mongo.Db).CreateCollection(ctx, c.cfg.MongoCollections.Products)
	if err != nil {
		if !utils.CheckErrMessages(err, constants.ErrMsgMongoCollectionAlreadyExists) {
			c.log.Warnf("(CreateCollection) err: {%v}", err)
		}
	}

	indexOptions := options.Index().SetSparse(true).SetUnique(true)
	index, err := mongoClient.Database(c.cfg.Mongo.Db).Collection(c.cfg.MongoCollections.Products).Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: consts.ProductIdIndex, Value: 1}},
		Options: indexOptions,
	})
	if err != nil && !utils.CheckErrMessages(err, constants.ErrMsgAlreadyExists) {
		c.log.Warnf("(CreateOne) err: {%v}", err)
	}
	c.log.Infof("(CreatedIndex) index: {%s}", index)

	list, err := mongoClient.Database(c.cfg.Mongo.Db).Collection(c.cfg.MongoCollections.Products).Indexes().List(ctx)
	if err != nil {
		c.log.Warnf("(initMongoDBCollections) [List] err: {%v}", err)
	}

	if list != nil {
		var results []bson.M
		if err := list.All(ctx, &results); err != nil {
			c.log.Warnf("(All) err: {%v}", err)
		}
		c.log.Infof("(indexes) results: {%#v}", results)
	}

	collections, err := mongoClient.Database(c.cfg.Mongo.Db).ListCollectionNames(ctx, bson.M{})
	if err != nil {
		c.log.Warnf("(ListCollections) err: {%v}", err)
	}
	c.log.Infof("(Collections) created collections: {%v}", collections)
}
