package main

import (
	"flag"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/constants"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/core"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger/zap"
	errorUtils "github.com/mehdihadeli/store-golang-microservice-sample/pkg/utils/error_utils"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/shared/server"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/shared/web"
	"log"
)

//https://github.com/swaggo/swag#how-to-use-it-with-gin

// @contact.name Mehdi Hadeli
// @contact.url https://github.com/mehdihadeli
// @title Catalogs Read-Service Api
// @version 1.0
// @description Catalogs Read-Service Api.
func main() {
	flag.Parse()

	//https://stackoverflow.com/questions/52103182/how-to-get-the-stacktrace-of-a-panic-and-store-as-a-variable
	defer errorUtils.HandlePanic()

	env := core.ConfigAppEnv(constants.Dev)

	cfg, err := config.InitConfig(env)
	if err != nil {
		log.Fatal(err)
	}

	logrusLogger := zap.NewZapLogger(cfg.Logger)
	logrusLogger.WithName(web.GetMicroserviceName(cfg))

	logrusLogger.Fatal(server.NewServer(logrusLogger, cfg).Run())
}
