package gomigrate

// https://github.com/golang-migrate/migrate/blob/856ea12df9d230b0145e23d951b7dbd6b86621cb/database/postgres/TUTORIAL.md#optional-run-migrations-within-your-go-app

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/migration"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/migration/contracts"

	"emperror.dev/errors"
	"github.com/golang-migrate/migrate/v4"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type goMigratePostgresMigrator struct {
	config     *migration.MigrationOptions
	db         *sql.DB
	datasource string
	logger     logger.Logger
	migration  *migrate.Migrate
}

func NewGoMigratorPostgres(
	config *migration.MigrationOptions,
	db *sql.DB,
	logger logger.Logger,
) (contracts.PostgresMigrationRunner, error) {
	if config.DBName == "" {
		return nil, errors.New("DBName is required in the config.")
	}

	datasource := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.DBName,
	)

	// In test environment, ewe need a fix for applying application working directory correctly. we will apply this in our environment setup process in `config/environment` file
	migration, err := migrate.New(fmt.Sprintf("file://%s", config.MigrationsDir), datasource)
	if err != nil {
		return nil, errors.WrapIf(err, "failed to initialize migrator")
	}

	return &goMigratePostgresMigrator{
		config:     config,
		db:         db,
		datasource: datasource,
		logger:     logger,
		migration:  migration,
	}, nil
}

func (m *goMigratePostgresMigrator) Up(_ context.Context, version uint) error {
	if m.config.SkipMigration {
		m.logger.Info("database migration skipped")

		return nil
	}

	err := m.executeCommand(migration.Up, version)

	if errors.Is(err, migrate.ErrNoChange) {
		return nil
	}

	if errors.Is(err, migrate.ErrNoChange) {
		return nil
	}

	if err != nil {
		return errors.WrapIf(err, "failed to migrate database")
	}

	m.logger.Info("migration finished")

	return nil
}

func (m *goMigratePostgresMigrator) Down(_ context.Context, version uint) error {
	if m.config.SkipMigration {
		m.logger.Info("database migration skipped")

		return nil
	}

	err := m.executeCommand(migration.Up, version)

	if errors.Is(err, migrate.ErrNoChange) {
		return nil
	}

	if err != nil {
		return errors.WrapIf(err, "failed to migrate database")
	}

	m.logger.Info("migration finished")

	return nil
}

func (m *goMigratePostgresMigrator) executeCommand(command migration.CommandType, version uint) error {
	var err error
	switch command {
	case migration.Up:
		if version == 0 {
			err = m.migration.Up()
		} else {
			err = m.migration.Migrate(version)
		}
	case migration.Down:
		if version == 0 {
			err = m.migration.Down()
		} else {
			err = m.migration.Migrate(version)
		}
	default:
		err = errors.New("invalid migration direction")
	}

	if err == migrate.ErrNoChange {
		return nil
	}
	if err != nil {
		return errors.WrapIf(err, "failed to migrate database")
	}

	return nil
}
