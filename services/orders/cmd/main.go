package main

import (
	"flag"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/constants"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/core"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger/zap"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/shared/server"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/shared/web"
	"log"
)

//https://github.com/swaggo/swag#how-to-use-it-with-gin

// @contact.name Mehdi Hadeli
// @contact.url https://github.com/mehdihadeli
// @title Orders Service Api
// @version 1.0
// @description Orders Service Api
func main() {
	flag.Parse()

	env := core.ConfigAppEnv(constants.Dev)

	cfg, err := config.InitConfig(env)
	if err != nil {
		log.Fatal(err)
	}

	appLogger := zap.NewZapLogger(cfg.Logger)
	appLogger.WithName(web.GetMicroserviceName(cfg))

	appLogger.Fatal(server.NewServer(appLogger, cfg).Run())
}
