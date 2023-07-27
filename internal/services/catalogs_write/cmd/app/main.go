package main

import (
	"flag"
	"log"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/constants"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/core"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger/zap"
	errorUtils "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/utils/error_utils"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/write_service/config"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/write_service/internal/shared/server"
)

// @contact.name Mehdi Hadeli
// @contact.url https://github.com/mehdihadeli
// @title Catalogs Write-Service Api
// @version 1.0
// @description Catalogs Write-Service Api.
func main() {
	flag.Parse()

	// https://stackoverflow.com/questions/52103182/how-to-get-the-stacktrace-of-a-panic-and-store-as-a-variable
	defer errorUtils.HandlePanic()

	env := core.ConfigAppEnv(constants.Dev)

	cfg, err := config.InitConfig(env)
	if err != nil {
		log.Fatal(err)
	}

	appLogger := zap.NewZapLogger(cfg.Logger)
	appLogger.WithName(cfg.GetMicroserviceName())
	appLogger.Fatal(server.NewServer(appLogger, cfg).Run())
}
