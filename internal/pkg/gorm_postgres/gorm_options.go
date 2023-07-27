package gormPostgres

import (
	"fmt"

	"github.com/iancoleman/strcase"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/config"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/config/environemnt"
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

func (h *GormOptions) Dns() string {
	datasource := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		h.User,
		h.Password,
		h.Host,
		h.Port,
		h.DBName,
	)

	return datasource
}

func provideConfig(environment environemnt.Environment) (*GormOptions, error) {
	return config.BindConfigKey[*GormOptions](optionName, environment)
}
