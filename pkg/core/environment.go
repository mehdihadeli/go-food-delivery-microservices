package core

import (
	"github.com/joho/godotenv"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/constants"
	"log"
	"os"
	"syscall"
)

func IsDevelopment() bool {
	env := os.Getenv("APP_ENV")
	return env == constants.Dev
}

func IsProduction() bool {
	env := os.Getenv("APP_ENV")
	return env == constants.Production
}

func GetEnvironment() string {
	if IsDevelopment() {
		return constants.Dev
	} else if IsProduction() {
		return constants.Production
	}
	
	return constants.Dev
}

func ConfigAppEnv(env string) string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	manualEnv := os.Getenv("APP_ENV")

	var envResult string
	if env == "" {
		envResult = constants.Dev
	} else {
		envResult = env
	}

	if manualEnv != "" {
		envResult = manualEnv
	}

	return envResult
}

func EnvString(key, fallback string) string {
	if value, ok := syscall.Getenv(key); ok {
		return value
	}
	return fallback
}
