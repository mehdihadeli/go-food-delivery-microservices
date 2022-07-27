package infrastructure

import (
	"context"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/mongodb"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
)

func (ic *infrastructureConfigurator) configMongo(ctx context.Context) (*mongo.Client, error, func()) {
	mongo, err := mongodb.NewMongoDB(ctx, ic.cfg.Mongo)
	if err != nil {
		return nil, errors.Wrap(err, "NewMongoDBConn"), nil
	}

	ic.log.Infof("(Mongo connected) SessionsInProgress: {%v}", mongo.MongoClient.NumberSessionsInProgress())

	return mongo.MongoClient, nil, func() {
		_ = mongo.Close() // nolint: errcheck
	}
}
