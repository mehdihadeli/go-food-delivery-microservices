package mongodb

import (
	"context"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/tracing"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

const (
	connectTimeout  = 30 * time.Second
	maxConnIdleTime = 3 * time.Minute
	minPoolSize     = 20
	maxPoolSize     = 300
)

type Config struct {
	URI      string `mapstructure:"uri"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Db       string `mapstructure:"db"`
	UseAuth  bool   `mapstructure:"useAuth"`
}

// NewMongoDBConn Create new MongoDB client
func NewMongoDBConn(ctx context.Context, cfg *Config) (*mongo.Client, error) {

	opt := options.Client().ApplyURI(cfg.URI).
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

	return client, nil
}

//https://stackoverflow.com/a/23650312/581476

func Paginate[T any](ctx context.Context, listQuery *utils.ListQuery, collection *mongo.Collection, filter interface{}) (*utils.ListResult[T], error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "mongodb.Paginate")

	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, errors.Wrap(err, "CountDocuments")
	}

	limit := int64(listQuery.GetLimit())
	skip := int64(listQuery.GetOffset())
	if filter == nil {
		filter = bson.D{}
	}
	cursor, err := collection.Find(ctx, filter, &options.FindOptions{
		Limit: &limit,
		Skip:  &skip,
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, errors.Wrap(err, "Find")
	}
	defer cursor.Close(ctx) // nolint: errcheck

	products := make([]*T, 0, listQuery.GetSize())

	for cursor.Next(ctx) {
		var prod T
		if err := cursor.Decode(&prod); err != nil {
			tracing.TraceErr(span, err)
			return nil, errors.Wrap(err, "Find")
		}
		products = append(products, &prod)
	}

	if err := cursor.Err(); err != nil {
		tracing.TraceErr(span, err)
		return nil, errors.Wrap(err, "cursor.Err")
	}

	return utils.NewListResult(products, listQuery.GetSize(), listQuery.GetPage(), count), nil
}
