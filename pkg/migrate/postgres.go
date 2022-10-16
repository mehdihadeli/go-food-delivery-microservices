package migrate

import (
	"context"
	"database/sql"
	"emperror.dev/errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
	"path/filepath"
	"runtime"
)

// Up executes all migrations found at the given source path against the
// database specified by given DSN.
func Up(config *MigrationConfig) error {
	if config.SkipMigration {
		zap.L().Info("database migration skipped")
		return nil
	}

	if config.DBName == "" {
		return errors.New("DBName is required in the config.")
	}

	err := createDB(config, context.Background())
	if err != nil {
		return err
	}

	datasource := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.DBName,
	)

	db, err := sql.Open("postgres", datasource)
	if err != nil {
		return err
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{
		MigrationsTable: config.VersionTable,
		DatabaseName:    config.DBName,
	})
	if err != nil {
		return fmt.Errorf("failed to initialize migrator: %w", err)
	}

	// determine the project's root path
	_, callerPath, _, _ := runtime.Caller(0) // nolint:dogsled

	// look for migrations source starting from project's root dir
	sourceURL := fmt.Sprintf(
		"file://%s/../../%s",
		filepath.ToSlash(filepath.Dir(callerPath)),
		filepath.ToSlash(config.MigrationsDir),
	)

	migration, err := migrate.NewWithDatabaseInstance(
		sourceURL,
		config.DBName,
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to initialize migrator: %w", err)
	}

	if config.TargetVersion == 0 {
		err = migration.Up()
	} else {
		err = migration.Migrate(config.TargetVersion)
	}

	if err == migrate.ErrNoChange {
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}
	zap.L().Info("migration finished")

	return nil
}

func createDB(cfg *MigrationConfig, ctx context.Context) error {
	// we should choose a default database in the connection, but because we don't have a database yet we specify postgres default database 'postgres'
	datasource := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		"postgres",
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
