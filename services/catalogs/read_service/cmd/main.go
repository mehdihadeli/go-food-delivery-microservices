package main

import (
	"flag"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/shared/server"
	"log"
	"os"
)

const dev = "development"
const production = "production"

// @contact.name Mehdi Hadeli
// @contact.url https://github.com/mehdihadeli
func main() {
	flag.Parse()

	env := os.Getenv("APP_ENV")
	if env == "" {
		env = dev
	}

	cfg, err := config.InitConfig(env)
	if err != nil {
		log.Fatal(err)
	}

	appLogger := logger.NewAppLogger(cfg.Logger)
	appLogger.InitLogger()
	appLogger.WithName("catalogs-services")

	appLogger.Fatal(server.NewServer(appLogger, cfg).Run())
}
