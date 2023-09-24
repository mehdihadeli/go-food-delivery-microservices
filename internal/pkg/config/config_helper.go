package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/config/environemnt"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/constants"
	typeMapper "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/reflection/type_mappper"

	"emperror.dev/errors"
	"github.com/caarlos0/env/v8"
	"github.com/mcuadros/go-defaults"
	"github.com/spf13/viper"
)

func BindConfig[T any](environments ...environemnt.Environment) (T, error) {
	return BindConfigKey[T]("", environments...)
}

func BindConfigKey[T any](
	configKey string,
	environments ...environemnt.Environment,
) (T, error) {
	var configPath string

	environment := environemnt.Environment("")
	if len(environments) > 0 {
		environment = environments[0]
	} else {
		environment = constants.Dev
	}

	cfg := typeMapper.GenericInstanceByT[T]()

	// this should set before reading config values from json file
	// https://github.com/mcuadros/go-defaults
	defaults.SetDefaults(cfg)

	// https://articles.wesionary.team/environment-variable-configuration-in-your-golang-project-using-viper-4e8289ef664d
	// when we `Set` a viper with string value, we should get it from viper with `viper.GetString`, elsewhere we get empty string
	// load `config path` from environment variable or viper internal registry
	configPathFromEnv := viper.GetString(constants.ConfigPath)

	if configPathFromEnv != "" {
		configPath = configPathFromEnv
	} else {
		// https://stackoverflow.com/questions/31873396/is-it-possible-to-get-the-current-root-of-package-structure-as-a-string-in-golan
		// https://stackoverflow.com/questions/18537257/how-to-get-the-directory-of-the-currently-running-file
		appRootPath := viper.GetString(constants.AppRootPath)

		d, err := searchForConfigFileDir(appRootPath, environment)
		if err != nil {
			return *new(T), err
		}

		configPath = d
	}

	// https://github.com/spf13/viper/issues/390#issuecomment-718756752
	viper.SetConfigName(fmt.Sprintf("config.%s", environment))
	viper.AddConfigPath(configPath)
	viper.SetConfigType(constants.Json)

	if err := viper.ReadInConfig(); err != nil {
		return *new(T), errors.WrapIf(err, "viper.ReadInConfig")
	}

	if len(configKey) == 0 {
		// load configs from config file to config object
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

	return cfg, nil
}

// searchForConfigFileDir searches for the first directory within the specified root directory and its subdirectories
// that contains a file named "config.%s.json" where "%s" is replaced with the provided environment string.
// It returns the path of the first directory that contains the config file or an error if no such directory is found.
//
// Parameters:
//
//	rootDir:      The root directory to start the search from.
//	environment:  The environment string to replace "%s" in the config file name.
//
// Returns:
//
//	string: The path of the directory containing the config file, or an empty string if not found.
//	error:  An error indicating any issues encountered during the search.
func searchForConfigFileDir(
	rootDir string,
	environment environemnt.Environment,
) (string, error) {
	var result string

	// Walk the directory tree starting from rootDir
	err := filepath.Walk(
		rootDir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// Check if the file is named "config.%s.json" (replace %s with the environment)
			if !info.IsDir() &&
				strings.EqualFold(
					info.Name(),
					fmt.Sprintf("config.%s.json", environment),
				) ||
				strings.EqualFold(
					info.Name(),
					fmt.Sprintf("config.%s.yaml", environment),
				) ||
				strings.EqualFold(
					info.Name(),
					fmt.Sprintf("config.%s.yml", environment),
				) {
				// Get the directory name containing the config file
				dir := filepath.Dir(path)
				result = dir
				return filepath.SkipDir // Skip further traversal
			}

			return nil
		},
	)

	if result != "" {
		return result, nil
	}

	return "", errors.WrapIf(err, "No directory with config file found")
}
