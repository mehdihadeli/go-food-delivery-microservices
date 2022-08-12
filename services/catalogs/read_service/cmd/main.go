package main

import (
	"flag"
	"github.com/joho/godotenv"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/constants"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger/logrous"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/shared/server"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/shared/web"
	"log"
	"os"
)

//https://github.com/swaggo/swag#how-to-use-it-with-gin

// @contact.name Mehdi Hadeli
// @contact.url https://github.com/mehdihadeli
// @title Catalogs Read-Service Api
// @version 1.0
// @description Catalogs Read-Service Api.
func main() {
	flag.Parse()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	env := os.Getenv("APP_ENV")
	if env == "" {
		env = constants.Dev
	}

	cfg, err := config.InitConfig(env)
	if err != nil {
		log.Fatal(err)
	}

	logrusLogger := logrous.NewLogrusLogger(cfg.Logger)
	logrusLogger.WithName(web.GetMicroserviceName(cfg))

	logrusLogger.Fatal(server.NewServer(logrusLogger, cfg).Run())
}
