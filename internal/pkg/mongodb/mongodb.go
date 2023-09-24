package mongodb

import (
	"context"
	"fmt"
	"time"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const (
	connectTimeout  = 60 * time.Second
	maxConnIdleTime = 3 * time.Minute
	minPoolSize     = 20
	maxPoolSize     = 300
)

// NewMongoDB Create new MongoDB client
func NewMongoDB(cfg *MongoDbOptions) (*mongo.Client, error) {
	uriAddress := fmt.Sprintf(
		"mongodb://%s:%s@%s:%d",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
	)
	opt := options.Client().ApplyURI(uriAddress).
		SetConnectTimeout(connectTimeout).
		SetMaxConnIdleTime(maxConnIdleTime).
		SetMinPoolSize(minPoolSize).
		SetMaxPoolSize(maxPoolSize)

	if cfg.UseAuth {
		opt = opt.SetAuth(
			options.Credential{Username: cfg.User, Password: cfg.Password},
		)
	}

	ctx := context.Background()
	client, err := mongo.Connect(ctx, opt)
	if err != nil {
		return nil, err
	}

	if cfg.EnableTracing {
		// add tracing
		opt.Monitor = otelmongo.NewMonitor()
	}

	// setup  https://github.com/Kamva/mgm
	err = mgm.SetDefaultConfig(nil, cfg.Database, opt)
	if err != nil {
		return nil, err
	}

	return client, nil
}
