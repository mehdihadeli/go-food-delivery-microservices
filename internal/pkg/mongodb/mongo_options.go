package mongodb

import (
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/config"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/config/environemnt"
	typeMapper "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/reflection/type_mappper"

	"github.com/iancoleman/strcase"
)

type MongoDbOptions struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
	UseAuth  bool   `mapstructure:"useAuth"`
}

func provideConfig(environment environemnt.Environment) (*MongoDbOptions, error) {
	optionName := strcase.ToLowerCamel(typeMapper.GetTypeNameByT[MongoDbOptions]())
	return config.BindConfigKey[*MongoDbOptions](optionName, environment)
}
