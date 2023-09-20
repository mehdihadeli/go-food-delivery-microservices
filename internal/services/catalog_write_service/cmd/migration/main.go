package main

import (
	"context"
	"os"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/config"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/config/environemnt"
	gormPostgres "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/gorm_postgres"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
	defaultLogger "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger/default_logger"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger/external/fxlog"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger/zap"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/migration"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/migration/contracts"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/migration/goose"
	appconfig "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/config"

	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

func init() {
	// Add flags to specify the version
	cmdUp.Flags().Uint("version", 0, "Migration version")
	cmdDown.Flags().Uint("version", 0, "Migration version")

	// Add commands to the root command
	rootCmd.AddCommand(cmdUp)
	rootCmd.AddCommand(cmdDown)
}

var (
	rootCmd = &cobra.Command{ //nolint:gochecknoglobals
		Use:   "migration",
		Short: "A tool for running migrations",
		Run: func(cmd *cobra.Command, args []string) {
			// Execute the "up" subcommand when no subcommand is specified
			if len(args) == 0 {
				cmd.SetArgs([]string{"up"})
				if err := cmd.Execute(); err != nil {
					defaultLogger.Logger.Error(err)
					os.Exit(1)
				}
			}
		},
	}

	cmdDown = &cobra.Command{ //nolint:gochecknoglobals
		Use:   "down",
		Short: "Run a down migration",
		Run: func(cmd *cobra.Command, args []string) {
			executeMigration(cmd, migration.Down)
		},
	}

	cmdUp = &cobra.Command{ //nolint:gochecknoglobals
		Use:   "up",
		Short: "Run an up migration",
		Run: func(cmd *cobra.Command, args []string) {
			executeMigration(cmd, migration.Up)
		},
	}
)

func executeMigration(cmd *cobra.Command, commandType migration.CommandType) {
	version, err := cmd.Flags().GetUint("version")
	if err != nil {
		defaultLogger.Logger.Fatal(err)
	}

	app := fx.New(
		config.ModuleFunc(environemnt.Development),
		zap.Module,
		fxlog.FxLogger,
		gormPostgres.Module,
		appconfig.Module,
		//// use go-migrate library for migration
		//gomigrate.Module,
		// use go-migrate library for migration
		goose.Module,
		fx.Invoke(func(migrationRunner contracts.PostgresMigrationRunner, logger logger.Logger) {
			logger.Info("Migration process started...")
			switch commandType {
			case migration.Up:
				err = migrationRunner.Up(context.Background(), version)
			case migration.Down:
				err = migrationRunner.Down(context.Background(), version)
			}
			if err != nil {
				logger.Fatalf("migration failed, err: %s", err)
			}
			logger.Info("Migration completed...")
		}),
	)

	err = app.Start(context.Background())
	if err != nil {
		defaultLogger.Logger.Fatal(err)
	}

	err = app.Stop(context.Background())
	if err != nil {
		defaultLogger.Logger.Fatal(err)
	}
}

func main() {
	defaultLogger.SetupDefaultLogger()

	if err := rootCmd.Execute(); err != nil {
		defaultLogger.Logger.Error(err)
		os.Exit(1)
	}
}
