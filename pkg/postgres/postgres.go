package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	DBName   string `yaml:"dbName"`
	SSLMode  bool   `yaml:"sslMode"`
	Password string `yaml:"password"`
}

const (
	maxConn           = 50
	healthCheckPeriod = 1 * time.Minute
	maxConnIdleTime   = 1 * time.Minute
	maxConnLifetime   = 3 * time.Minute
	minConns          = 10
	lazyConnect       = false
)

// NewPgxConn func for connection to PostgreSQL database.
func NewPgxConn(cfg *Config) (*pgxpool.Pool, error) {
	ctx := context.Background()
	var dataSourceName string

	if cfg.DBName == "" {
		dataSourceName = fmt.Sprintf("host=%s port=%s user=%s password=%s",
			cfg.Host,
			cfg.Port,
			cfg.User,
			cfg.Password,
		)
	} else {
		dataSourceName = fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s",
			cfg.Host,
			cfg.Port,
			cfg.User,
			cfg.DBName,
			cfg.Password,
		)
	}

	poolCfg, err := pgxpool.ParseConfig(dataSourceName)
	if err != nil {
		return nil, err
	}

	poolCfg.MaxConns = maxConn
	poolCfg.HealthCheckPeriod = healthCheckPeriod
	poolCfg.MaxConnIdleTime = maxConnIdleTime
	poolCfg.MaxConnLifetime = maxConnLifetime
	poolCfg.MinConns = minConns
	poolCfg.LazyConnect = lazyConnect

	connPool, err := pgxpool.ConnectConfig(ctx, poolCfg)
	if err != nil {
		return nil, errors.Wrap(err, "pgx.ConnectConfig")
	}

	return connPool, nil
}

// NewSqlxConn func for connection to PostgreSQL database.
func NewSqlxConn(cfg *Config) (*sqlx.DB, error) {
	// Define database connection settings.
	maxConn, _ := strconv.Atoi(os.Getenv("DB_MAX_CONNECTIONS"))
	maxIdleConn, _ := strconv.Atoi(os.Getenv("DB_MAX_IDLE_CONNECTIONS"))
	maxLifetimeConn, _ := strconv.Atoi(os.Getenv("DB_MAX_LIFETIME_CONNECTIONS"))

	var dataSourceName string

	if cfg.DBName == "" {
		dataSourceName = fmt.Sprintf("host=%s port=%s user=%s password=%s",
			cfg.Host,
			cfg.Port,
			cfg.User,
			cfg.Password,
		)
	} else {
		dataSourceName = fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s",
			cfg.Host,
			cfg.Port,
			cfg.User,
			cfg.DBName,
			cfg.Password,
		)
	}

	// Define database connection for PostgreSQL.
	db, err := sqlx.Connect("pgx", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("error, not connected to database, %w", err)
	}

	// Set database connection settings.
	db.SetMaxOpenConns(maxConn)                           // the default is 0 (unlimited)
	db.SetMaxIdleConns(maxIdleConn)                       // defaultMaxIdleConns = 2
	db.SetConnMaxLifetime(time.Duration(maxLifetimeConn)) // 0, connections are reused forever

	// Try to ping database.
	if err := db.Ping(); err != nil {
		defer db.Close() // close database connection
		return nil, fmt.Errorf("error, not sent ping to database, %w", err)
	}

	return db, nil
}
