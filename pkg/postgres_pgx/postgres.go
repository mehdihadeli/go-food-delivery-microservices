package postgres

import (
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/migrations"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"time"
)

//Ref:https://github.com/henvic/pgxtutorial
// https://aiven.io/blog/aiven-for-postgresql-for-your-go-application
import (
	"context"
	"github.com/jackc/pgx/v4"
)

type Config struct {
	Host       string     `mapstructure:"host"`
	Port       string     `mapstructure:"port"`
	User       string     `mapstructure:"user"`
	DBName     string     `mapstructure:"dbName"`
	SSLMode    bool       `mapstructure:"sslMode"`
	Password   string     `mapstructure:"password"`
	Migrations Migrations `mapstructure:"migrations"`
}

type Migrations struct {
	MigrationsDirectory string `mapstructure:"migrationsDir"`
	VersionTable        string `mapstructure:"versionTable"`
	SchemaVersion       uint   `mapstructure:"schemaVersion"`
	SkipMigration       bool   `mapstructure:"skipMigration"`
}

const (
	maxConn           = 50
	healthCheckPeriod = 1 * time.Minute
	maxConnIdleTime   = 1 * time.Minute
	maxConnLifetime   = 3 * time.Minute
	minConns          = 10
	lazyConnect       = false
)

// txCtx key.
type txCtx struct{}

// connCtx key.
type connCtx struct{}

type Pgx struct {
	ConnPool        *pgxpool.Pool
	DB              *sql.DB
	SquirrelBuilder squirrel.StatementBuilderType
	GoquBuilder     *goqu.SelectDataset
}

// NewPgxPoolConn func for connection to PostgreSQL database.
func NewPgxPoolConn(cfg *Config, logger pgx.Logger, logLevel pgx.LogLevel) (*Pgx, error) {
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

	// https://henvic.dev/posts/go-postgres/
	// https://aiven.io/blog/aiven-for-postgresql-for-your-go-application
	// https://jbrandhorst.com/post/postgres/

	if logger != nil {
		poolCfg.ConnConfig.Logger = logger
	}

	if logLevel != 0 {
		poolCfg.ConnConfig.LogLevel = logLevel
	}

	connPool, err := pgxpool.ConnectConfig(ctx, poolCfg)
	if err != nil {
		return nil, errors.Wrap(err, "pgx.ConnectConfig")
	}

	pgxConfig, err := pgx.ParseConfig(dataSourceName)
	if err != nil {
		return nil, err
	}

	// https://github.com/jackc/pgx/issues/737#issuecomment-640075332
	//db, err := sql.Open("pgx", dataSourceName)
	db := stdlib.OpenDB(*pgxConfig) // db.Conn().Raw() - get a connection from the pool with stdlib and sql/database

	// goqu
	dialect := goqu.Dialect("postgres")
	database := dialect.DB(db)
	goquBuilder := database.From()

	// squirrel
	squirrelBuilder := squirrel.StatementBuilder.
		PlaceholderFormat(squirrel.Dollar).RunWith(db)

	p := &Pgx{ConnPool: connPool, DB: db, SquirrelBuilder: squirrelBuilder, GoquBuilder: goquBuilder}

	if cfg.Migrations.SkipMigration {
		zap.L().Info("database migration skipped")
		return p, nil
	}

	mp := migrations.MigrationParams{
		DbName:        cfg.DBName,
		VersionTable:  cfg.Migrations.VersionTable,
		MigrationsDir: cfg.Migrations.MigrationsDirectory,
		TargetVersion: cfg.Migrations.SchemaVersion,
	}

	if err = migrations.RunMigration(db, mp); err != nil {
		return nil, err
	}

	return p, nil
}

func (db *Pgx) Close() {
	db.ConnPool.Close()
	_ = db.DB.Close()
}

// conn returns a PostgreSQL transaction if one exists.
// If not, returns a connection if a connection has been acquired by calling WithAcquire.
// Otherwise, it returns *pgxpool.Pool which acquires the connection and closes it immediately after a SQL command is executed.
func (db *Pgx) conn(ctx context.Context) PGXQuerier {
	if tx, ok := ctx.Value(txCtx{}).(pgx.Tx); ok && tx != nil {
		return tx
	}
	if res, ok := ctx.Value(connCtx{}).(*pgxpool.Conn); ok && res != nil {
		return res
	}
	return db.ConnPool
}

// TransactionContext returns a copy of the parent context which begins a transaction
// to PostgreSQL.
//
// Once the transaction is over, you must call db.Commit(ctx) to make the changes effective.
// This might live in the go-pkg/postgres package later for the sake of code reuse.
func (db *Pgx) TransactionContext(ctx context.Context) (context.Context, error) {
	tx, err := db.conn(ctx).Begin(ctx)
	if err != nil {
		return nil, err
	}
	return context.WithValue(ctx, txCtx{}, tx), nil
}

// Commit transaction from context.
func (db *Pgx) Commit(ctx context.Context) error {
	if tx, ok := ctx.Value(txCtx{}).(pgx.Tx); ok && tx != nil {
		return tx.Commit(ctx)
	}
	return errors.New("context has no transaction")
}

// Rollback transaction from context.
func (db *Pgx) Rollback(ctx context.Context) error {
	if tx, ok := ctx.Value(txCtx{}).(pgx.Tx); ok && tx != nil {
		return tx.Rollback(ctx)
	}
	return errors.New("context has no transaction")
}

// WithAcquire returns a copy of the parent context which acquires a connection
// to PostgreSQL from pgxpool to make sure commands executed in series reuse the
// same database connection.
//
// To release the connection back to the pool, you must call postgres.Release(ctx).
//
// Example:
// dbCtx := db.WithAcquire(ctx)
// defer postgres.Release(dbCtx)
func (db *Pgx) WithAcquire(ctx context.Context) (dbCtx context.Context, err error) {
	if _, ok := ctx.Value(connCtx{}).(*pgxpool.Conn); ok {
		panic("context already has a connection acquired")
	}
	res, err := db.ConnPool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	return context.WithValue(ctx, connCtx{}, res), nil
}

// Release PostgreSQL connection acquired by context back to the pool.
func (db *Pgx) Release(ctx context.Context) {
	if res, ok := ctx.Value(connCtx{}).(*pgxpool.Conn); ok && res != nil {
		res.Release()
	}
}
