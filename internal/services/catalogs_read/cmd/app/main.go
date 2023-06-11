package main

import (
	"flag"
	"os"

	"github.com/spf13/cobra"
	"go.uber.org/fx"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/config"
	defaultLogger "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger/default_logger"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger/external/fxlog"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger/zap"
	errorUtils "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/utils/error_utils"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/read_service/internal/app"
)

const version = "2.0.0"

var rootCmd = &cobra.Command{
	Use:     "ecommerce-microservices",
	Version: version,
	Short:   "ecommerce-microservices",
	Run: func(cmd *cobra.Command, args []string) {
		// https://stackoverflow.com/questions/52103182/how-to-get-the-stacktrace-of-a-panic-and-store-as-a-variable
		fxApp := fx.New( // setup fxlog logger
			zap.Module,
			config.Module,
			fxlog.FxLogger,

			app.Module)

		fxApp.Run()
	},
}

// https://github.com/swaggo/swag#how-to-use-it-with-gin

// @contact.name Mehdi Hadeli
// @contact.url https://github.com/mehdihadeli
// @title Catalogs Read-Service Api
// @version 1.0
// @description Catalogs Read-Service Api.
func main() {
	flag.Parse()
	defer errorUtils.HandlePanic()

	if err := rootCmd.Execute(); err != nil {
		defaultLogger.Logger.Fatal(err)
		os.Exit(1)
	}

	// appLogger := zap.NewZapLogger(cfg.Logger)
	// appLogger.WithName(cfg.GetMicroserviceName())
	// appLogger.Fatal(server.NewServer(appLogger, cfg).Run())
}
