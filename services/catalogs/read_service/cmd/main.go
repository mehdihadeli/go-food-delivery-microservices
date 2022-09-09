package main

import (
	"flag"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/core"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger/zap"
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

	env := core.ConfigAppEnv()

	cfg, err := config.InitConfig(env)
	if err != nil {
		log.Fatal(err)
	}

	logrusLogger := zap.NewZapLogger(cfg.Logger)
	logrusLogger.WithName(web.GetMicroserviceName(cfg))

	logrusLogger.Fatal(server.NewServer(logrusLogger, cfg).Run())
}
