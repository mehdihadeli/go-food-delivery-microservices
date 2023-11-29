package postgresgorm

import (
	"database/sql"
	"fmt"

	defaultlogger "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger/defaultlogger"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger/external/gromlog"

	"emperror.dev/errors"
	"github.com/glebarez/sqlite"
	gormPostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/opentelemetry/tracing"
)

func NewGorm(cfg *GormOptions) (*gorm.DB, error) {
	if cfg.DBName == "" {
		return nil, errors.New("DBName is required in the config.")
	}

	err := createDB(cfg)
	if err != nil {
		return nil, err
	}

	if cfg.UseInMemory {
		db, err := createInMemoryDB()

		return db, err
	}

	dataSourceName := fmt.Sprintf(
		"host=%s port=%d user=%s dbname=%s password=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.DBName,
		cfg.Password,
	)

	gormDb, err := gorm.Open(
		gormPostgres.Open(dataSourceName),
		&gorm.Config{
			Logger: gromlog.NewGormCustomLogger(defaultlogger.GetLogger()),
		},
	)
	if err != nil {
		return nil, err
	}

	// add tracing to gorm
	if cfg.EnableTracing {
		err = gormDb.Use(tracing.NewPlugin())
	}

	return gormDb, nil
}

func createInMemoryDB() (*gorm.DB, error) {
	// https://gorm.io/docs/connecting_to_the_database.html#SQLite
	// https://github.com/glebarez/sqlite
	// https://www.connectionstrings.com/sqlite/
	db, err := gorm.Open(
		sqlite.Open(":memory:"),
		&gorm.Config{
			Logger: gromlog.NewGormCustomLogger(defaultlogger.GetLogger()),
		})

	return db, err
}

func NewSQLDB(orm *gorm.DB) (*sql.DB, error) { return orm.DB() }

func createDB(cfg *GormOptions) error {
	var db *sql.DB

	// we should choose a default database in the connection, but because we don't have a database yet we specify postgres default database 'postgres'
	dataSourceName := fmt.Sprintf(
		"host=%s port=%d user=%s dbname=%s password=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		"postgres",
		cfg.Password,
	)
	postgresGormDB, err := gorm.Open(
		gormPostgres.Open(dataSourceName),
		&gorm.Config{
			Logger: gromlog.NewGormCustomLogger(defaultlogger.GetLogger()),
		},
	)
	if err != nil {
		return err
	}

	db, err = postgresGormDB.DB()

	if err != nil {
		return err
	}

	rows, err := db.Query(
		fmt.Sprintf(
			"SELECT 1 FROM  pg_catalog.pg_database WHERE datname='%s'",
			cfg.DBName,
		),
	)
	if err != nil {
		return err
	}

	var exists int
	if rows.Next() {
		err = rows.Scan(&exists)
		if err != nil {
			return err
		}
	}

	if exists == 1 {
		return nil
	}

	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", cfg.DBName))
	if err != nil {
		return err
	}

	defer db.Close()

	return nil
}
