package core

import (
	"github.com/joho/godotenv"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/constants"
	"log"
	"os"
)

func IsDevelopment() bool {
	env := os.Getenv("APP_ENV")
	return env == constants.Dev
}

func IsProduction() bool {
	env := os.Getenv("APP_ENV")
	return env == constants.Production
}

func ConfigAppEnv() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	env := os.Getenv("APP_ENV")
	if env == "" {
		env = constants.Dev
	}

	return env
}