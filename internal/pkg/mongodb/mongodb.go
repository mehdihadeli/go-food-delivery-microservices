package mongodb

import (
	"context"
	"fmt"
	"time"

	"emperror.dev/errors"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/otel/tracing"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/utils"

	_ "github.com/golang-migrate/migrate/v4/source/file"

	"go.mongodb.org/mongo-driver/mongo"
)

const (
	connectTimeout  = 30 * time.Second
	maxConnIdleTime = 3 * time.Minute
	minPoolSize     = 20
	maxPoolSize     = 300
)

type MongoDbConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
	UseAuth  bool   `mapstructure:"useAuth"`
}

// NewMongoDB Create new MongoDB client
func NewMongoDB(ctx context.Context, cfg *MongoDbConfig) (*mongo.Client, error) {
	uriAddres := fmt.Sprintf("mongodb://%s:%s@%s:%d", cfg.User, cfg.Password, cfg.Host, cfg.Port)
	opt := options.Client().ApplyURI(uriAddres).
		SetConnectTimeout(connectTimeout).
		SetMaxConnIdleTime(maxConnIdleTime).
		SetMinPoolSize(minPoolSize).
		SetMaxPoolSize(maxPoolSize)

	if cfg.UseAuth {
		opt = opt.SetAuth(options.Credential{Username: cfg.User, Password: cfg.Password})
	}

	client, err := mongo.NewClient(opt)
	if err != nil {
		return nil, err
	}

	if err := client.Connect(ctx); err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	// setup  https://github.com/Kamva/mgm
	err = mgm.SetDefaultConfig(nil, cfg.Database, opt)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// https://stackoverflow.com/a/23650312/581476

func Paginate[T any](ctx context.Context, listQuery *utils.ListQuery, collection *mongo.Collection, filter interface{}) (*utils.ListResult[T], error) {
	ctx, span := tracing.Tracer.Start(ctx, "mongodb.Paginate")
	defer span.End()

	if filter == nil {
		filter = bson.D{}
	}

	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, tracing.TraceErrFromSpan(span, errors.WrapIf(err, "CountDocuments"))
	}

	limit := int64(listQuery.GetLimit())
	skip := int64(listQuery.GetOffset())

	cursor, err := collection.Find(ctx, filter, &options.FindOptions{
		Limit: &limit,
		Skip:  &skip,
	})
	if err != nil {
		return nil, tracing.TraceErrFromSpan(span, errors.WrapIf(err, "Find"))
	}
	defer cursor.Close(ctx) // nolint: errcheck

	products := make([]T, 0, listQuery.GetSize())

	for cursor.Next(ctx) {
		var prod T
		if err := cursor.Decode(&prod); err != nil {
			return nil, tracing.TraceErrFromSpan(span, errors.WrapIf(err, "Find"))
		}
		products = append(products, prod)
	}

	if err := cursor.Err(); err != nil {
		return nil, tracing.TraceErrFromSpan(span, errors.WrapIf(err, "cursor.Err"))
	}

	return utils.NewListResult[T](products, listQuery.GetSize(), listQuery.GetPage(), count), nil
}
