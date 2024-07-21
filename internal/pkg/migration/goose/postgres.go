package goose

// https://github.com/pressly/goose#embedded-sql-migrations

import (
	"context"
	"database/sql"
	"errors"
	"strconv"

	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/logger"
	migration "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/migration"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/migration/contracts"

	"github.com/pressly/goose/v3"
)

type goosePostgresMigrator struct {
	config *migration.MigrationOptions
	db     *sql.DB
	logger logger.Logger
}

func NewGoosePostgres(
	config *migration.MigrationOptions,
	db *sql.DB,
	logger logger.Logger,
) contracts.PostgresMigrationRunner {
	goose.SetBaseFS(nil)

	return &goosePostgresMigrator{config: config, db: db, logger: logger}
}

func (m *goosePostgresMigrator) Up(_ context.Context, version uint) error {
	err := m.executeCommand(migration.Up, version)

	return err
}

func (m *goosePostgresMigrator) Down(_ context.Context, version uint) error {
	err := m.executeCommand(migration.Down, version)

	return err
}

func (m *goosePostgresMigrator) executeCommand(
	command migration.CommandType,
	version uint,
) error {
	switch command {
	case migration.Up:
		if version == 0 {
			// In test environment, we need a fix for applying application working directory correctly. we will apply this in our environment setup process in `config/environment` file
			return goose.Run("up", m.db, m.config.MigrationsDir)
		}

		return goose.Run(
			"up-to VERSION ",
			m.db,
			m.config.MigrationsDir,
			strconv.FormatUint(uint64(version), 10),
		)
	case migration.Down:
		if version == 0 {
			return goose.Run("down", m.db, m.config.MigrationsDir)
		}

		return goose.Run(
			"down-to VERSION ",
			m.db,
			m.config.MigrationsDir,
			strconv.FormatUint(uint64(version), 10),
		)
	default:
		return errors.New("invalid migration direction")
	}
}
