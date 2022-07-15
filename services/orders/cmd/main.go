package main

import (
	"flag"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"log"
	"os"
	"thub.com/mehdihadeli/store-golang-microservice-sample/services/orders/config"
	"thub.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/shared/server"
	"thub.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/shared/web"
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
	appLogger.WithName(web.GetMicroserviceName(cfg))

	appLogger.Fatal(server.NewServer(appLogger, cfg).Run())
}
