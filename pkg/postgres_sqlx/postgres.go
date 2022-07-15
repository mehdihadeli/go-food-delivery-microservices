package postgres_sqlx

import (
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/doug-martin/goqu/v9"
	_ "github.com/jackc/pgx/v4/stdlib" // load pgx driver for PostgreSQL
	"github.com/jmoiron/sqlx"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/migrations"
	"go.uber.org/zap"
	"os"
	"strconv"
	"time"
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

	// stdlib package doesn't have a compat layer for pgxpool
	// so had to use standard sql api for pool configuration.
	db.SetMaxOpenConns(maxConn)                           // the default is 0 (unlimited)
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

func (db *Sqlx) Migrate() error {
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

	if err := migrations.RunMigration(db.DB, mp); err != nil {
		return err
	}

	return nil
}

func (db *Sqlx) Close() {
	_ = db.DB.Close()
	_ = db.SqlxDB.Close()
}
