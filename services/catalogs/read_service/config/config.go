package config

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/grpc"
	customEcho "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/custom_echo"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/otel"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/config"
	"os"
	"path/filepath"
	"runtime"

	"emperror.dev/errors"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/constants"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/elasticsearch"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/eventstroredb"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/mongodb"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/probes"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/redis"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/tracing"
	"github.com/spf13/viper"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config", "", "catalogs read_position microservice config path")
}

type Config struct {
	DeliveryType     string                         `mapstructure:"deliveryType" env:"DeliveryType"`
	ServiceName      string                         `mapstructure:"serviceName" env:"ServiceName"`
	Logger           *logger.LogConfig              `mapstructure:"logger" envPrefix:"Logger_"`
	GRPC             *grpc.GrpcConfig               `mapstructure:"grpc" envPrefix:"GRPC_"`
	Http             *customEcho.EchoHttpConfig     `mapstructure:"http" envPrefix:"Http_"`
	Context          Context                        `mapstructure:"context" envPrefix:"Context_"`
	Redis            *redis.Config                  `mapstructure:"redis" envPrefix:"Redis_"`
	RabbitMQ         *config.RabbitMQConfig         `mapstructure:"rabbitmq" envPrefix:"RabbitMQ_"`
	Probes           probes.Config                  `mapstructure:"probes" envPrefix:"Probes_"`
	Jaeger           *tracing.Config                `mapstructure:"jaeger" envPrefix:"Jaeger_"`
	OTel             *otel.OpenTelemetryConfig      `mapstructure:"otel" envPrefix:"OTel_"`
	EventStoreConfig eventstroredb.EventStoreConfig `mapstructure:"eventStoreConfig" envPrefix:"EventStoreConfig_"`
	Elastic          elasticsearch.Config           `mapstructure:"elastic" envPrefix:"Elastic_"`
	ElasticIndexes   ElasticIndexes                 `mapstructure:"elasticIndexes" envPrefix:"ElasticIndexes_"`
	Mongo            *mongodb.MongoDbConfig         `mapstructure:"mongo" envPrefix:"Mongo_"`
	MongoCollections MongoCollections               `mapstructure:"mongoCollections" envPrefix:"MongoCollections_"`
}

type Context struct {
	Timeout int `mapstructure:"timeout" env:"Timeout"`
}

type MongoCollections struct {
	Products string `mapstructure:"products" validate:"required" env:"Products"`
}

type ElasticIndexes struct {
	Products string `mapstructure:"products" validate:"required" env:"Products"`
}

func InitConfig(environment string) (*Config, error) {
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

	//https://github.com/spf13/viper/issues/390#issuecomment-718756752
	viper.SetConfigName(fmt.Sprintf("config.%s", environment))
	viper.AddConfigPath(configPath)
	viper.SetConfigType(constants.Json)

	if err := viper.ReadInConfig(); err != nil {
		return nil, errors.WrapIf(err, "viper.ReadInConfig")
	}

	if err := viper.Unmarshal(cfg); err != nil {
		return nil, errors.WrapIf(err, "viper.Unmarshal")
	}

	if err := env.Parse(cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}

	grpcPort := os.Getenv(constants.GrpcPort)
	if grpcPort != "" {
		cfg.GRPC.Port = grpcPort
	}

	mongoURI := os.Getenv(constants.MongoDbURI)
	if mongoURI != "" {
		//cfg.Mongo.URI = "mongodb://host.docker.internal:27017"
		cfg.Mongo.URI = mongoURI
	}

	redisAddr := os.Getenv(constants.RedisAddr)
	if redisAddr != "" {
		cfg.Redis.Addr = redisAddr
	}

	jaegerAddr := os.Getenv(constants.JaegerHostPort)
	if jaegerAddr != "" {
		cfg.Jaeger.HostPort = jaegerAddr
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
