package mongodb

import (
	"context"
	"emperror.dev/errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mongodb"
	"github.com/kamva/mgm/v3"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/core/data"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/otel/tracing"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"time"

	_ "github.com/golang-migrate/migrate/v4/source/file"

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
	config      *MongoDbConfig
}

type MongoDbConfig struct {
	Host       string               `mapstructure:"host"`
	Port       int                  `mapstructure:"port"`
	User       string               `mapstructure:"user"`
	Password   string               `mapstructure:"password"`
	Database   string               `mapstructure:"database"`
	UseAuth    bool                 `mapstructure:"useAuth"`
	Migrations data.MigrationParams `mapstructure:"migrations"`
}

// NewMongoDB Create new MongoDB client
func NewMongoDB(ctx context.Context, cfg *MongoDbConfig) (*MongoDb, error) {
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

	mp := data.MigrationParams{
		DbName:        m.config.Database,
		VersionTable:  m.config.Migrations.VersionTable,
		MigrationsDir: m.config.Migrations.MigrationsDir,
		TargetVersion: m.config.Migrations.TargetVersion,
	}

	d, err := mongodb.WithInstance(m.MongoClient, &mongodb.Config{DatabaseName: mp.DbName, MigrationsCollection: mp.VersionTable})
	if err != nil {
		return fmt.Errorf("failed to initialize migrator: %w", err)
	}

	mig, err := migrate.NewWithDatabaseInstance("file://"+mp.MigrationsDir, mp.DbName, d)
	if err != nil {
		return fmt.Errorf("failed to initialize migrator: %w", err)
	}

	if mp.TargetVersion == 0 {
		err = mig.Up()
	} else {
		err = mig.Migrate(mp.TargetVersion)
	}

	if err == migrate.ErrNoChange {
		return nil
	}

	zap.L().Info("migration finished")
	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}
	return nil
}

//https://stackoverflow.com/a/23650312/581476

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
