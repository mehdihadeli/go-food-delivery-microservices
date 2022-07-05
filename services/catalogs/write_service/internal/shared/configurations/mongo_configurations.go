package configurations

import (
	"context"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/mongodb"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/utils"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/consts"
	catalog_constants "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/shared/constants"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (ic *infrastructureConfigurator) configMongo(ctx context.Context) (*mongo.Client, error, func()) {
	mongoClient, err := mongodb.NewMongoDBConn(ctx, ic.cfg.Mongo)
	if err != nil {
		return nil, errors.Wrap(err, "NewMongoDBConn"), nil
	}

	ic.log.Infof("(Mongo connected) SessionsInProgress: {%v}", mongoClient.NumberSessionsInProgress())

	ic.initMongoDBCollections(ctx, mongoClient)

	return mongoClient, nil, func() {
		_ = mongoClient.Disconnect(ctx) // nolint: errcheck
	}
}

func (ic *infrastructureConfigurator) initMongoDBCollections(ctx context.Context, mongoClient *mongo.Client) {
	err := mongoClient.Database(ic.cfg.Mongo.Db).CreateCollection(ctx, ic.cfg.MongoCollections.Products)
	if err != nil {
		if !utils.CheckErrMessages(err, catalog_constants.ErrMsgMongoCollectionAlreadyExists) {
			ic.log.Warnf("(CreateCollection) err: {%v}", err)
		}
	}

	indexOptions := options.Index().SetSparse(true).SetUnique(true)
	index, err := mongoClient.Database(ic.cfg.Mongo.Db).Collection(ic.cfg.MongoCollections.Products).Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: consts.ProductIdIndex, Value: 1}},
		Options: indexOptions,
	})
	if err != nil && !utils.CheckErrMessages(err, catalog_constants.ErrMsgAlreadyExists) {
		ic.log.Warnf("(CreateOne) err: {%v}", err)
	}
	ic.log.Infof("(CreatedIndex) index: {%s}", index)

	list, err := mongoClient.Database(ic.cfg.Mongo.Db).Collection(ic.cfg.MongoCollections.Products).Indexes().List(ctx)
	if err != nil {
		ic.log.Warnf("(initMongoDBCollections) [List] err: {%v}", err)
	}

	if list != nil {
		var results []bson.M
		if err := list.All(ctx, &results); err != nil {
			ic.log.Warnf("(All) err: {%v}", err)
		}
		ic.log.Infof("(indexes) results: {%#v}", results)
	}

	collections, err := mongoClient.Database(ic.cfg.Mongo.Db).ListCollectionNames(ctx, bson.M{})
	if err != nil {
		ic.log.Warnf("(ListCollections) err: {%v}", err)
	}
	ic.log.Infof("(Collections) created collections: {%v}", collections)
}
