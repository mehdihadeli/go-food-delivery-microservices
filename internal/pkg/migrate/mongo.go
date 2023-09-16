package migrate

import (
	"context"
	"fmt"
	"path/filepath"
	"runtime"

	mongodb2 "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/mongodb"

	"emperror.dev/errors"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mongodb"
	"go.uber.org/zap"
)

func (config *MigrationConfig) Migrate(ctx context.Context) error {
	if config.SkipMigration {
		zap.L().Info("database migration skipped")
		return nil
	}

	if config.DBName == "" {
		return errors.New("DBName is required in the config.")
	}

	db, err := mongodb2.NewMongoDB(&mongodb2.MongoDbOptions{
		Host:     config.Host,
		Port:     config.Port,
		User:     config.User,
		Password: config.Password,
		Database: config.DBName,
		UseAuth:  false,
	})
	if err != nil {
		return err
	}

	driver, err := mongodb.WithInstance(
		db,
		&mongodb.Config{DatabaseName: config.DBName, MigrationsCollection: config.VersionTable},
	)
	if err != nil {
		return fmt.Errorf("failed to initialize migrator: %w", err)
	}

	// determine the project's root path
	_, callerPath, _, _ := runtime.Caller(1) // nolint:dogsled

	// look for migrations source starting from project's root dir
	sourceURL := fmt.Sprintf(
		"file://%s/../../%s",
		filepath.ToSlash(filepath.Dir(callerPath)),
		filepath.ToSlash(config.MigrationsDir),
	)

	mig, err := migrate.NewWithDatabaseInstance(sourceURL, config.DBName, driver)
	if err != nil {
		return fmt.Errorf("failed to initialize migrator: %w", err)
	}

	if config.TargetVersion == 0 {
		err = mig.Up()
	} else {
		err = mig.Migrate(config.TargetVersion)
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
