package gormPostgres

import (
	"github.com/iancoleman/strcase"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/config"
	typeMapper "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/reflection/type_mappper"
)

var optionName = strcase.ToLowerCamel(typeMapper.GetTypeNameByT[GormOptions]())

type GormOptions struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	DBName   string `mapstructure:"dbName"`
	SSLMode  bool   `mapstructure:"sslMode"`
	Password string `mapstructure:"password"`
}

func provideConfig() (*GormOptions, error) {
	return config.BindConfigKey[*GormOptions](optionName)
}
