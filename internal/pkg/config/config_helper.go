package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/config/environment"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/constants"
	typeMapper "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/reflection/typemapper"

	"emperror.dev/errors"
	"github.com/caarlos0/env/v8"
	"github.com/mcuadros/go-defaults"
	"github.com/spf13/viper"
)

func BindConfig[T any](environments ...environment.Environment) (T, error) {
	return BindConfigKey[T]("", environments...)
}

func BindConfigKey[T any](
	configKey string,
	environments ...environment.Environment,
) (T, error) {
	var configPath string

	currentEnv := environment.Environment("")
	if len(environments) > 0 {
		currentEnv = environments[0]
	} else {
		currentEnv = constants.Dev
	}

	cfg := typeMapper.GenericInstanceByT[T]()

	// this should set before reading config values from json file
	// https://github.com/mcuadros/go-defaults
	defaults.SetDefaults(cfg)

	// https://articles.wesionary.team/environment-variable-configuration-in-your-golang-project-using-viper-4e8289ef664d
	// when we `Set` a viper with string value, we should get it from viper with `viper.GetString`, elsewhere we get empty string
	// load `config path` from env variable or viper internal registry
	configPathFromEnv := viper.GetString(constants.ConfigPath)

	if configPathFromEnv != "" {
		configPath = configPathFromEnv
	} else {
		// https://stackoverflow.com/questions/31873396/is-it-possible-to-get-the-current-root-of-package-structure-as-a-string-in-golan
		// https://stackoverflow.com/questions/18537257/how-to-get-the-directory-of-the-currently-running-file
		appRootPath := viper.GetString(constants.AppRootPath)
		if appRootPath == "" {
			appRootPath = environment.GetProjectRootWorkingDirectory()
		}

		d, err := searchForConfigFileDir(appRootPath, currentEnv)
		if err != nil {
			return *new(T), err
		}

		configPath = d
	}

	// https://github.com/spf13/viper/issues/390#issuecomment-718756752
	viper.SetConfigName(fmt.Sprintf("config.%s", currentEnv))
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
	env environment.Environment,
) (string, error) {
	var result string

	// Walk the directory tree starting from rootDir
	err := filepath.Walk(
		rootDir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// Check if the file is named "config.%s.json" (replace %s with the env)
			if !info.IsDir() &&
				strings.EqualFold(
					info.Name(),
					fmt.Sprintf("config.%s.json", env),
				) ||
				strings.EqualFold(
					info.Name(),
					fmt.Sprintf("config.%s.yaml", env),
				) ||
				strings.EqualFold(
					info.Name(),
					fmt.Sprintf("config.%s.yml", env),
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
