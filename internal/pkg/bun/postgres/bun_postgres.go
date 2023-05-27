package postgres

import (
    "database/sql"
    "fmt"

    "emperror.dev/errors"
    "github.com/uptrace/bun"
    "github.com/uptrace/bun/dialect/pgdialect"
    "github.com/uptrace/bun/driver/pgdriver"

    bun2 "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/bun"
    // loading bun's official Postgres driver.
    _ "github.com/uptrace/bun/driver/pgdriver"
)

func NewBunDB(cfg *bun2.BunConfig) (*bun.DB, error) {
	if cfg.DBName == "" {
		return nil, errors.New("DBName is required in the config.")
	}

	err := createDB(cfg)
	if err != nil {
		return nil, err
	}

	//https://bun.uptrace.dev/postgres/#pgdriver
	datasource := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
	)
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(datasource)))

	//pgconn := pgdriver.NewConnector(
	//	pgdriver.WithNetwork("tcp"),
	//	pgdriver.WithAddr("localhost:5437"),
	//	pgdriver.WithTLSConfig(&tls.Config{InsecureSkipVerify: true}),
	//	pgdriver.WithUser("test"),
	//	pgdriver.WithPassword("test"),
	//	pgdriver.WithDatabase("test"),
	//	pgdriver.WithApplicationName("myapp"),
	//	pgdriver.WithTimeout(5*time.Second),
	//	pgdriver.WithDialTimeout(5*time.Second),
	//	pgdriver.WithReadTimeout(5*time.Second),
	//	pgdriver.WithWriteTimeout(5*time.Second),
	//	pgdriver.WithConnParams(map[string]interface{}{
	//		"search_path": "my_search_path",
	//	}),
	//)
	//sqldb := sql.OpenDB(pgconn)

	db := bun.NewDB(sqldb, pgdialect.New())

	return db, nil
}

func createDB(cfg *bun2.BunConfig) error {
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
	rows, err := sqldb.Query(fmt.Sprintf("SELECT 1 FROM  pg_catalog.pg_database WHERE datname='%s'", cfg.DBName))
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
