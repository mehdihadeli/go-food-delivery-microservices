package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"go.uber.org/fx"

	customEcho "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/custom_echo"
	defaultLogger "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger/default_logger"
	errorUtils "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/utils/error_utils"
	application "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/write_service/internal/shared/app"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/write_service/internal/shared/configurations/catalogs"
)

const version = "1.0.0"

var rootCmd = &cobra.Command{
	Version: version,
	Use:     "catalogs-write-service",
	Short:   "catalogs-write-service",
	Run: func(cmd *cobra.Command, args []string) {
		// configure dependencies
		appBuilder := application.NewCatalogsWriteApplicationBuilder()
		appBuilder.ProvideModule(catalogs.Module)

		app := appBuilder.Build()

		app.ResolveFunc(func(echo customEcho.EchoHttpServer) {
			fmt.Print(echo)
		})

		app.RegisterHook(func(lifecycle fx.Lifecycle) {
			lifecycle.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					return nil
				},
				OnStop: func(ctx context.Context) error {
					// some cleanup if exists
					return nil
				},
			})
		})

		// configure application
		app.ConfigureCatalogs()

		app.MapCatalogsEndpoints()

		app.Run()
	},
}

// https://github.com/swaggo/swag#how-to-use-it-with-gin

// @contact.name Mehdi Hadeli
// @contact.url https://github.com/mehdihadeli
// @title Catalogs Write-Service Api
// @version 1.0
// @description Catalogs Write-Service Api.
func main() {
	flag.Parse()
	defer errorUtils.HandlePanic()

	if err := rootCmd.Execute(); err != nil {
		defaultLogger.Logger.Fatal(err)
		os.Exit(1)
	}
}
