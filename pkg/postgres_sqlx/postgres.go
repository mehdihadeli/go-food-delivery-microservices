package postgres_sqlx

import (
	"context"
	"database/sql"
	"emperror.dev/errors"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/doug-martin/goqu/v9"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/jackc/pgx/v4/stdlib" // load pgx driver for PostgreSQL
	"github.com/jmoiron/sqlx"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/core/data"
	"go.uber.org/zap"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Host       string               `mapstructure:"host"`
	Port       string               `mapstructure:"port"`
	User       string               `mapstructure:"user"`
	DBName     string               `mapstructure:"dbName"`
	SSLMode    bool                 `mapstructure:"sslMode"`
	Password   string               `mapstructure:"password"`
	Migrations data.MigrationParams `mapstructure:"migrations"`
}

type Sqlx struct {
	SqlxDB          *sqlx.DB
	DB              *sql.DB
	SquirrelBuilder squirrel.StatementBuilderType
	GoquBuilder     *goqu.SelectDataset
	config          *Config
}

// NewSqlxConn creates a database connection with appropriate pool configuration
// and runs migration to prepare database.
//
// Migration will be omitted if appropriate config parameter set.
func NewSqlxConn(cfg *Config) (*Sqlx, error) {
	// Define database connection settings.
	maxConn, _ := strconv.Atoi(os.Getenv("DB_MAX_CONNECTIONS"))
	maxIdleConn, _ := strconv.Atoi(os.Getenv("DB_MAX_IDLE_CONNECTIONS"))
	maxLifetimeConn, _ := strconv.Atoi(os.Getenv("DB_MAX_LIFETIME_CONNECTIONS"))

	var dataSourceName string

	if cfg.DBName == "" {
		return nil, errors.New("DBName is required in the config.")
	}

	err := createDB(cfg)

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

	// Define database connection for PostgreSQL.
	db, err := sqlx.Connect("pgx", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("error, not connected to database, %w", err)
	}

	// stdlib package doesn't have a compat layer for pgxpool
	// so had to use standard sql api for pool configuration.
	db.SetMaxOpenConns(maxConn)                           // the defaultLogger is 0 (unlimited)
	db.SetMaxIdleConns(maxIdleConn)                       // defaultMaxIdleConns = 2
	db.SetConnMaxLifetime(time.Duration(maxLifetimeConn)) // 0, connections are reused forever

	// Try to ping database.
	if err := db.Ping(); err != nil {
		defer db.Close() // close database connection
		return nil, fmt.Errorf("error, not sent ping to database, %w", err)
	}

	//squirrel
	squirrelBuilder := squirrel.StatementBuilder.
		PlaceholderFormat(squirrel.Dollar).RunWith(db)

	// goqu
	dialect := goqu.Dialect("postgres")
	database := dialect.DB(db)
	goquBuilder := database.From()

	sqlx := &Sqlx{DB: db.DB, SqlxDB: db, SquirrelBuilder: squirrelBuilder, GoquBuilder: goquBuilder, config: cfg}

	return sqlx, nil
}

func createDB(cfg *Config) error {
	datasource := fmt.Sprintf("host=%s port=%s user=%s password=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
	)
	db, err := sqlx.Connect("pgx", datasource)
	if err != nil {
		return fmt.Errorf("error, not connected to database, %w", err)
	}

	var exists int
	rows, err := db.Query(fmt.Sprintf("SELECT 1 FROM  pg_catalog.pg_database WHERE datname='%s'", cfg.DBName))
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

	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", cfg.DBName))
	if err != nil {
		return err
	}

	defer db.Close()

	return nil
}

func (db *Sqlx) Migrate() error {
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

	if err := runPostgresMigration(db.DB, mp); err != nil {
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

func (db *Sqlx) Close() {
	_ = db.DB.Close()
	_ = db.SqlxDB.Close()
}

// Ref:https://dev.to/techschoolguru/a-clean-way-to-implement-database-transaction-in-golang-2ba

// ExecTx executes a transaction with provided function.
func (db *Sqlx) ExecTx(ctx context.Context, fn func(*Sqlx) error) error {
	tx, err := db.DB.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return err
	}

	err = fn(db)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}
