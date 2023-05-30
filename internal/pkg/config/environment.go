package config

import (
	"log"
	"os"
	"syscall"

	"github.com/joho/godotenv"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/constants"
)

type Environment string

func ConfigAppEnv(environments ...string) Environment {
	environment := ""
	if len(environments) > 0 {
		environment = environments[0]
	} else {
		environment = constants.Dev
	}

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	manualEnv := os.Getenv("APP_ENV")

	if manualEnv != "" {
		environment = manualEnv
	}

	return Environment(environment)
}

func (env Environment) IsDevelopment() bool {
	return env == constants.Dev
}

func (env Environment) IsProduction() bool {
	return env == constants.Production
}

func (env Environment) GetEnvironmentName() string {
	return string(env)
}

func EnvString(key, fallback string) string {
	if value, ok := syscall.Getenv(key); ok {
		return value
	}
	return fallback
}
