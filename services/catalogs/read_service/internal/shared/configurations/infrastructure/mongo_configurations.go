package infrastructure

import (
	"context"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/mongodb"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
)

func (ic *infrastructureConfigurator) configMongo(ctx context.Context) (*mongo.Client, error, func()) {
	mongoClient, err := mongodb.NewMongoDBConn(ctx, ic.cfg.Mongo)
	if err != nil {
		return nil, errors.Wrap(err, "NewMongoDBConn"), nil
	}

	ic.log.Infof("(Mongo connected) SessionsInProgress: {%v}", mongoClient.NumberSessionsInProgress())

	return mongoClient, nil, func() {
		_ = mongoClient.Disconnect(ctx) // nolint: errcheck
	}
}
