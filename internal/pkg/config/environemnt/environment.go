package environemnt

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/constants"

	"emperror.dev/errors"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Environment string

var (
	Development = Environment(constants.Dev)
	Test        = Environment(constants.Test)
	Production  = Environment(constants.Production)
)

func ConfigAppEnv(environments ...Environment) Environment {
	environment := Environment("")
	if len(environments) > 0 {
		environment = environments[0]
	} else {
		environment = Development
	}

	// setup viper to read from os environment with `viper.Get`
	viper.AutomaticEnv()

	// https://articles.wesionary.team/environment-variable-configuration-in-your-golang-project-using-viper-4e8289ef664d
	// load environment variables form .env files to system environment variables, it just find `.env` file in our current `executing working directory` in our app for example `catalogs_service`
	err := loadEnvFilesRecursive()
	if err != nil {
		log.Printf(".env file cannot be found, err: %v", err)
	}

	setRootWorkingDirectoryEnvironment()

	manualEnv := os.Getenv(constants.AppEnv)

	if manualEnv != "" {
		environment = Environment(manualEnv)
	}

	return environment
}

func (env Environment) IsDevelopment() bool {
	return env == Development
}

func (env Environment) IsProduction() bool {
	return env == Production
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

func loadEnvFilesRecursive() error {
	// Start from the current working directory
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	// Keep traversing up the directory hierarchy until you find an ".env" file
	for {
		envFilePath := filepath.Join(dir, ".env")
		err := godotenv.Load(envFilePath)

		if err == nil {
			// .env file found and loaded
			return nil
		}

		// Move up one directory level
		parentDir := filepath.Dir(dir)
		if parentDir == dir {
			// Reached the root directory, stop searching
			break
		}

		dir = parentDir
	}

	return errors.New(".env file not found in the project hierarchy")
}

func setRootWorkingDirectoryEnvironment() {
	// https://articles.wesionary.team/environment-variable-configuration-in-your-golang-project-using-viper-4e8289ef664d
	// when we `Set` a viper with string value, we should get it from viper with `viper.GetString`, elsewhere we get empty string
	// viper will get it from `os env` or a .env file
	pn := viper.GetString(constants.PROJECT_NAME_ENV)
	if pn == "" {
		log.Fatalf(
			"can't find environment variable `%s` in the system environment variables or a .env file",
			constants.PROJECT_NAME_ENV,
		)
	}

	// set root working directory of our app in the viper
	// https://stackoverflow.com/a/47785436/581476
	wd, _ := os.Getwd()

	for !strings.HasSuffix(wd, pn) {
		wd = filepath.Dir(wd)
	}

	absoluteRootWorkingDirectory, _ := filepath.Abs(wd)

	// when we `Set` a viper with string value, we should get it from viper with `viper.GetString`, elsewhere we get empty string
	viper.Set(constants.AppRootPath, absoluteRootWorkingDirectory)

	configPath := viper.GetString(constants.ConfigPath)

	// check for existing `CONFIG_PATH` variable in system environment variables - we can set it directly in .env file or system environments
	if configPath == "" {
		configPath := filepath.Join(absoluteRootWorkingDirectory, "config")

		// when we `Set` a viper with string value, we should get it from viper with `viper.GetString`, elsewhere we get empty string
		viper.Set(constants.ConfigPath, configPath)
	}
}
