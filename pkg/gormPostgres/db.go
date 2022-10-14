package gormPostgres

import (
	"context"
	"database/sql"
	"emperror.dev/errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/core/data"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/otel/tracing"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/utils"
	"go.uber.org/zap"
	gorm_postgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"strings"
)

type Config struct {
	Host       string               `mapstructure:"host"`
	Port       int                  `mapstructure:"port"`
	User       string               `mapstructure:"user"`
	DBName     string               `mapstructure:"dbName"`
	SSLMode    bool                 `mapstructure:"sslMode"`
	Password   string               `mapstructure:"password"`
	Migrations data.MigrationParams `mapstructure:"migrations"`
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

	dataSourceName = fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s",
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
	// we should choose a default database in the connection, but because we don't have a database yet we specify postgres default database 'postgres'
	datasource := fmt.Sprintf("postgres://%s:%s@%s:%d/postgres?sslmode=disable",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
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

	mp := data.MigrationParams{
		DbName:        db.config.DBName,
		VersionTable:  db.config.Migrations.VersionTable,
		MigrationsDir: db.config.Migrations.MigrationsDir,
		TargetVersion: db.config.Migrations.TargetVersion,
	}

	d, _ := db.DB.DB()
	if err := runPostgresMigration(d, mp); err != nil {
		return err
	}

	return nil
}

func runPostgresMigration(db *sql.DB, p data.MigrationParams) error {
	d, err := postgres.WithInstance(db, &postgres.Config{
		MigrationsTable: p.VersionTable,
		DatabaseName:    p.DbName,
	})
	if err != nil {
		return fmt.Errorf("failed to initialize migrator: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://"+p.MigrationsDir, p.DbName, d)
	if err != nil {
		return fmt.Errorf("failed to initialize migrator: %w", err)
	}

	if p.TargetVersion == 0 {
		err = m.Up()
	} else {
		err = m.Migrate(p.TargetVersion)
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

//Ref: https://dev.to/rafaelgfirmino/pagination-using-gorm-scopes-3k5f

func Paginate[T any](ctx context.Context, listQuery *utils.ListQuery, db *gorm.DB) (*utils.ListResult[T], error) {
	ctx, span := tracing.Tracer.Start(ctx, "gorm.Paginate")
	defer span.End()

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
		return nil, tracing.TraceErrFromSpan(span, errors.WrapIf(err, "error in finding products."))
	}

	return utils.NewListResult[T](items, listQuery.GetSize(), listQuery.GetPage(), totalRows), nil
}
