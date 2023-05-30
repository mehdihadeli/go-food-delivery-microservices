package config

import (
	"fmt"
	"os"
	"path/filepath"

	"emperror.dev/errors"
	"github.com/caarlos0/env/v6"
	"github.com/spf13/viper"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/constants"
	typeMapper "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/reflection/type_mappper"
)

func BindConfig[T any](environments ...string) (T, error) {
	return BindConfigKey[T]("", environments...)
}

func BindConfigKey[T any](configKey string, environments ...string) (T, error) {
	var configPath string

	environment := ""
	if len(environments) > 0 {
		environment = environments[0]
	} else {
		environment = constants.Dev
	}

	configPathFromEnv := os.Getenv(constants.ConfigPath)
	if configPathFromEnv != "" {
		configPath = configPathFromEnv
	} else {
		// https://stackoverflow.com/questions/31873396/is-it-possible-to-get-the-current-root-of-package-structure-as-a-string-in-golan
		// https://stackoverflow.com/questions/18537257/how-to-get-the-directory-of-the-currently-running-file
		d, err := getConfigRootPath()
		if err != nil {
			return *new(T), err
		}

		configPath = d
	}

	cfg := typeMapper.GenericInstanceByT[T]()

	// https://github.com/spf13/viper/issues/390#issuecomment-718756752
	viper.SetConfigName(fmt.Sprintf("config.%s", environment))
	viper.AddConfigPath(configPath)
	viper.SetConfigType(constants.Json)

	if err := viper.ReadInConfig(); err != nil {
		return *new(T), errors.WrapIf(err, "viper.ReadInConfig")
	}

	if len(configKey) == 0 {
		if err := viper.Unmarshal(cfg); err != nil {
			return *new(T), errors.WrapIf(err, "viper.Unmarshal")
		}
	} else {
		if err := viper.UnmarshalKey(configKey, cfg); err != nil {
			return *new(T), errors.WrapIf(err, "viper.Unmarshal")
		}
	}

	if err := env.Parse(cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}

	return cfg, nil
}

func getConfigRootPath() (string, error) {
	// Get the current working directory
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Get the absolute path of the executed project directory
	projectPath, err := filepath.Abs(wd)
	if err != nil {
		return "", err
	}

	// Get the path to the "config" folder within the project directory
	configPath := filepath.Join(projectPath, "config")

	return configPath, nil
}
