package gormPostgres

import (
	"context"
	"emperror.dev/errors"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/migrations"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/tracing"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/utils"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	gorm_postgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"strings"
)

type Config struct {
	Host       string                     `mapstructure:"host"`
	Port       string                     `mapstructure:"port"`
	User       string                     `mapstructure:"user"`
	DBName     string                     `mapstructure:"dbName"`
	SSLMode    bool                       `mapstructure:"sslMode"`
	Password   string                     `mapstructure:"password"`
	Migrations migrations.MigrationParams `mapstructure:"migrations"`
}

type Gorm struct {
	DB     *gorm.DB
	config *Config
}

func NewGorm(cfg *Config) (*Gorm, error) {

	var dataSourceName string
	ctx := context.Background()

	if cfg.DBName == "" {
		return nil, errors.New("DBName is required in the config.")
	}

	err := createDB(cfg, ctx)

	if err != nil {
		return nil, err
	}

	dataSourceName = fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.DBName,
		cfg.Password,
	)

	gormDb, err := gorm.Open(gorm_postgres.Open(dataSourceName), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	return &Gorm{DB: gormDb, config: cfg}, nil
}

func (db *Gorm) Close() {
	d, _ := db.DB.DB()
	_ = d.Close()
}

func createDB(cfg *Config, ctx context.Context) error {
	datasource := fmt.Sprintf("host=%s port=%s user=%s password=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
	)

	poolCfg, err := pgxpool.ParseConfig(datasource)
	if err != nil {
		return err
	}

	connPool, err := pgxpool.ConnectConfig(ctx, poolCfg)
	if err != nil {
		return errors.WrapIf(err, "pgx.ConnectConfig")
	}

	var exists int
	rows, err := connPool.Query(context.Background(), fmt.Sprintf("SELECT 1 FROM  pg_catalog.pg_database WHERE datname='%s'", cfg.DBName))
	if err != nil {
		return err
	}

	if rows.Next() {
		err = rows.Scan(&exists)
		if err != nil {
			return err
		}
	}

	if exists == 1 {
		return nil
	}

	_, err = connPool.Exec(context.Background(), fmt.Sprintf("CREATE DATABASE %s", cfg.DBName))
	if err != nil {
		return err
	}

	defer connPool.Close()

	return nil
}

func (db *Gorm) Migrate() error {
	if db.config.Migrations.SkipMigration {
		zap.L().Info("database migration skipped")
		return nil
	}

	mp := migrations.MigrationParams{
		DbName:        db.config.DBName,
		VersionTable:  db.config.Migrations.VersionTable,
		MigrationsDir: db.config.Migrations.MigrationsDir,
		TargetVersion: db.config.Migrations.TargetVersion,
	}

	d, _ := db.DB.DB()
	if err := migrations.RunPostgresMigration(d, mp); err != nil {
		return err
	}

	return nil
}

//Ref: https://dev.to/rafaelgfirmino/pagination-using-gorm-scopes-3k5f

func Paginate[T any](ctx context.Context, listQuery *utils.ListQuery, db *gorm.DB) (*utils.ListResult[T], error) {

	span, ctx := opentracing.StartSpanFromContext(ctx, "gorm.Paginate")

	var items []T
	var totalRows int64
	db.Model(items).Count(&totalRows)

	// generate where query
	query := db.Offset(listQuery.GetOffset()).Limit(listQuery.GetLimit()).Order(listQuery.GetOrderBy())

	if listQuery.Filters != nil {
		for _, filter := range listQuery.Filters {
			column := filter.Field
			action := filter.Comparison
			value := filter.Value

			switch action {
			case "equals":
				whereQuery := fmt.Sprintf("%s = ?", column)
				query = query.Where(whereQuery, value)
				break
			case "contains":
				whereQuery := fmt.Sprintf("%s LIKE ?", column)
				query = query.Where(whereQuery, "%"+value+"%")
				break
			case "in":
				whereQuery := fmt.Sprintf("%s IN (?)", column)
				queryArray := strings.Split(value, ",")
				query = query.Where(whereQuery, queryArray)
				break

			}
		}
	}

	if err := query.Find(&items).Error; err != nil {
		tracing.TraceErr(span, err)
		return nil, errors.WrapIf(err, "error in finding products.")
	}

	return utils.NewListResult[T](items, listQuery.GetSize(), listQuery.GetPage(), totalRows), nil
}
