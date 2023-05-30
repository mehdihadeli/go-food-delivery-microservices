package main

import (
	"flag"

	"go.uber.org/fx"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/config"
	customEcho "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/custom_echo"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger/zap"
	errorUtils "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/utils/error_utils"
	appconfig "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/read_service/config"
)

// https://github.com/swaggo/swag#how-to-use-it-with-gin

// @contact.name Mehdi Hadeli
// @contact.url https://github.com/mehdihadeli
// @title Catalogs Read-Service Api
// @version 1.0
// @description Catalogs Read-Service Api.
func main() {
	flag.Parse()

	// https://stackoverflow.com/questions/52103182/how-to-get-the-stacktrace-of-a-panic-and-store-as-a-variable
	defer errorUtils.HandlePanic()

	app := fx.New(
		// infrastructure setup
		config.Module,
		zap.Module,

		// application setup
		customEcho.Module,
		appconfig.Module,
	)

	app.Run()

	//cfg, err := serviceconfig.InitConfig(env)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//appLogger := zap.NewZapLogger(cfg.Logger)
	//appLogger.WithName(cfg.GetMicroserviceName())
	//appLogger.Fatal(server.NewServer(appLogger, cfg).Run())
}
