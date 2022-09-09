package infrastructure

import (
	"context"
	"emperror.dev/errors"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
)

func (ic *infrastructureConfigurator) configMongo(ctx context.Context) (*mongo.Client, error, func()) {
	mongo, err := mongodb.NewMongoDB(ctx, ic.cfg.Mongo)
	if err != nil {
		return nil, errors.WrapIf(err, "NewMongoDBConn"), nil
	}

	ic.log.Infof("(Mongo connected) SessionsInProgress: {%v}", mongo.MongoClient.NumberSessionsInProgress())

	return mongo.MongoClient, nil, func() {
		_ = mongo.Close() // nolint: errcheck
	}
}
