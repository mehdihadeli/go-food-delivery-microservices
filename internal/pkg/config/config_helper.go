package config

import (
	"fmt"
	"os"
	"path/filepath"

	"emperror.dev/errors"
	"github.com/caarlos0/env/v8"
	"github.com/mcuadros/go-defaults"
	"github.com/spf13/viper"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/config/environemnt"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/constants"
	typeMapper "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/reflection/type_mappper"
)

func BindConfig[T any](environments ...environemnt.Environment) (T, error) {
	return BindConfigKey[T]("", environments...)
}

func BindConfigKey[T any](configKey string, environments ...environemnt.Environment) (T, error) {
	var configPath string

	environment := environemnt.Environment("")
	if len(environments) > 0 {
		environment = environments[0]
	} else {
		environment = constants.Dev
	}

	// https://articles.wesionary.team/environment-variable-configuration-in-your-golang-project-using-viper-4e8289ef664d
	configPathFromEnv := viper.Get(constants.ConfigPath)
	if configPathFromEnv != nil {
		configPath = configPathFromEnv.(string)
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

	viper.AutomaticEnv()

	// https://github.com/caarlos0/env
	if err := env.Parse(cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}

	// https://github.com/mcuadros/go-defaults
	defaults.SetDefaults(cfg)

	return cfg, nil
}

func getConfigRootPath() (string, error) {
	// Get the current working directory
	// Getwd gives us the current working directory that we are running our app with `go run` command. in goland we can specify this working directory for the project
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	fmt.Println(fmt.Sprintf("Current working directory is: %s", wd))

	// Get the absolute path of the executed project directory
	absCurrentDir, err := filepath.Abs(wd)
	if err != nil {
		return "", err
	}

	// Get the path to the "config" folder within the project directory
	configPath := filepath.Join(absCurrentDir, "config")

	return configPath, nil
}
