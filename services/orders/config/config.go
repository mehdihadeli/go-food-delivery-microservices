package config

import (
	"emperror.dev/errors"
	"flag"
	"fmt"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/constants"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/eventstroredb"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/grpc"
	customEcho "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/custom_echo"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/mongodb"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/otel"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/probes"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/config"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"runtime"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config", "", "catalogs write microservice config path")
}

type Config struct {
	DeliveryType     string                          `mapstructure:"deliveryType"`
	ServiceName      string                          `mapstructure:"serviceName"`
	Logger           *logger.LogConfig               `mapstructure:"logger"`
	GRPC             *grpc.GrpcConfig                `mapstructure:"grpc"`
	Http             *customEcho.EchoHttpConfig      `mapstructure:"http"`
	Context          Context                         `mapstructure:"context"`
	Probes           probes.Config                   `mapstructure:"probes"`
	OTel             *otel.OpenTelemetryConfig       `mapstructure:"otel" envPrefix:"OTel_"`
	RabbitMQ         *config.RabbitMQConfig          `mapstructure:"rabbitmq" envPrefix:"RabbitMQ_"`
	EventStoreConfig *eventstroredb.EventStoreConfig `mapstructure:"eventStoreConfig"`
	Subscriptions    *Subscriptions                  `mapstructure:"subscriptions"`
	Mongo            *mongodb.MongoDbConfig          `mapstructure:"mongo" envPrefix:"Mongo_"`
	MongoCollections MongoCollections                `mapstructure:"mongoCollections" envPrefix:"MongoCollections_"`
}

type Context struct {
	Timeout int `mapstructure:"timeout"`
}

type MongoCollections struct {
	Orders string `mapstructure:"orders" validate:"required" env:"Orders"`
}

type Subscriptions struct {
	OrderSubscription *Subscription `mapstructure:"orderSubscription"`
}

type Subscription struct {
	Prefix         []string `mapstructure:"prefix" validate:"required"`
	SubscriptionId string   `mapstructure:"subscriptionId" validate:"required"`
}

func InitConfig(env string) (*Config, error) {
	if configPath == "" {
		configPathFromEnv := os.Getenv(constants.ConfigPath)
		if configPathFromEnv != "" {
			configPath = configPathFromEnv
		} else {
			//https://stackoverflow.com/questions/31873396/is-it-possible-to-get-the-current-root-of-package-structure-as-a-string-in-golan
			//https://stackoverflow.com/questions/18537257/how-to-get-the-directory-of-the-currently-running-file
			d, err := dirname()
			if err != nil {
				return nil, err
			}

			configPath = d
		}
	}

	cfg := &Config{}

	viper.SetConfigName(fmt.Sprintf("config.%s", env))
	viper.AddConfigPath(configPath)
	viper.SetConfigType(constants.Yaml)

	if err := viper.ReadInConfig(); err != nil {
		return nil, errors.WrapIf(err, "viper.ReadInConfig")
	}

	if err := viper.Unmarshal(cfg); err != nil {
		return nil, errors.WrapIf(err, "viper.Unmarshal")
	}

	grpcPort := os.Getenv(constants.GrpcPort)
	if grpcPort != "" {
		cfg.GRPC.Port = grpcPort
	}

	jaegerPort := os.Getenv(constants.JaegerPort)
	if jaegerPort != "" {
		cfg.OTel.JaegerExporterConfig.AgentPort = jaegerPort
	}

	jaegerHost := os.Getenv(constants.JaegerHost)
	if jaegerHost != "" {
		cfg.OTel.JaegerExporterConfig.AgentHost = jaegerHost
	}

	return cfg, nil
}

func filename() (string, error) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", errors.New("unable to get the current filename")
	}
	return filename, nil
}

func dirname() (string, error) {
	filename, err := filename()
	if err != nil {
		return "", err
	}
	return filepath.Dir(filename), nil
}
