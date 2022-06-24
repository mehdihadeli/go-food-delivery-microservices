package main

import (
	"flag"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/server"
	"log"
)

// @contact.name Mehdi Hadeli
// @contact.url https://github.com/mehdihadeli
func main() {
	flag.Parse()

	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatal(err)
	}

	appLogger := logger.NewAppLogger(cfg.Logger)
	appLogger.InitLogger()
	appLogger.WithName("catalogs-services")

	appLogger.Fatal(server.NewServer(appLogger, cfg).Run())
}
