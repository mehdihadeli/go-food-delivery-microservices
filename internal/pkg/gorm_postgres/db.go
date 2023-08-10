package gormPostgres

import (
	"database/sql"
	"fmt"

	"emperror.dev/errors"
	"github.com/uptrace/bun/driver/pgdriver"
	gormPostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"

	defaultLogger "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger/default_logger"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger/external/gromlog"
)

func NewGorm(cfg *GormOptions) (*gorm.DB, error) {
	if cfg.DBName == "" {
		return nil, errors.New("DBName is required in the config.")
	}

	err := createDB(cfg)
	if err != nil {
		return nil, err
	}

	var dataSourceName string
	dataSourceName = fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.DBName,
		cfg.Password,
	)

	gormDb, err := gorm.Open(
		gormPostgres.Open(dataSourceName),
		&gorm.Config{Logger: gromlog.NewGormCustomLogger(defaultLogger.Logger)},
	)
	if err != nil {
		return nil, err
	}

	return gormDb, nil
}

func createDB(cfg *GormOptions) error {
	// we should choose a default database in the connection, but because we don't have a database yet we specify postgres default database 'postgres'
	datasource := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		"postgres",
	)

	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(datasource)))

	var exists int
	rows, err := sqldb.Query(
		fmt.Sprintf("SELECT 1 FROM  pg_catalog.pg_database WHERE datname='%s'", cfg.DBName),
	)
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

	_, err = sqldb.Exec(fmt.Sprintf("CREATE DATABASE %s", cfg.DBName))
	if err != nil {
		return err
	}

	defer sqldb.Close()

	return nil
}
