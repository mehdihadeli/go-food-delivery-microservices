package mongodb

import (
	"context"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/migrations"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/tracing"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

const (
	connectTimeout  = 30 * time.Second
	maxConnIdleTime = 3 * time.Minute
	minPoolSize     = 20
	maxPoolSize     = 300
)

type MongoDb struct {
	MongoClient *mongo.Client
	config      *Config
}

type Config struct {
	URI        string                     `mapstructure:"uri"`
	User       string                     `mapstructure:"user"`
	Password   string                     `mapstructure:"password"`
	Db         string                     `mapstructure:"db"`
	UseAuth    bool                       `mapstructure:"useAuth"`
	Migrations migrations.MigrationParams `mapstructure:"migrations"`
}

// NewMongoDB Create new MongoDB client
func NewMongoDB(ctx context.Context, cfg *Config) (*MongoDb, error) {

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

	return &MongoDb{MongoClient: client}, nil
}

func (m *MongoDb) Close() error {
	return m.MongoClient.Disconnect(context.Background())
}

func (m *MongoDb) Migrate() error {
	if m.config.Migrations.SkipMigration {
		zap.L().Info("database migration skipped")
		return nil
	}

	mp := migrations.MigrationParams{
		DbName:        m.config.Db,
		VersionTable:  m.config.Migrations.VersionTable,
		MigrationsDir: m.config.Migrations.MigrationsDir,
		TargetVersion: m.config.Migrations.TargetVersion,
	}

	if err := migrations.RunMongoMigration(m.MongoClient, mp); err != nil {
		return err
	}

	return nil
}

//https://stackoverflow.com/a/23650312/581476

func Paginate[T any](ctx context.Context, listQuery *utils.ListQuery, collection *mongo.Collection, filter interface{}) (*utils.ListResult[T], error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "mongodb.Paginate")

	if filter == nil {
		filter = bson.D{}
	}

	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, errors.Wrap(err, "CountDocuments")
	}

	limit := int64(listQuery.GetLimit())
	skip := int64(listQuery.GetOffset())

	cursor, err := collection.Find(ctx, filter, &options.FindOptions{
		Limit: &limit,
		Skip:  &skip,
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, errors.Wrap(err, "Find")
	}
	defer cursor.Close(ctx) // nolint: errcheck

	products := make([]T, 0, listQuery.GetSize())

	for cursor.Next(ctx) {
		var prod T
		if err := cursor.Decode(&prod); err != nil {
			tracing.TraceErr(span, err)
			return nil, errors.Wrap(err, "Find")
		}
		products = append(products, prod)
	}

	if err := cursor.Err(); err != nil {
		tracing.TraceErr(span, err)
		return nil, errors.Wrap(err, "cursor.Err")
	}

	return utils.NewListResult[T](products, listQuery.GetSize(), listQuery.GetPage(), count), nil
}
