package config

import (
	"fmt"
	"net/url"

	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/config"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/config/environment"
	typeMapper "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/reflection/typemapper"

	"github.com/iancoleman/strcase"
)

var optionName = strcase.ToLowerCamel(typeMapper.GetGenericTypeNameByT[EchoHttpOptions]())

type EchoHttpOptions struct {
	Port                string   `mapstructure:"port"                validate:"required" env:"TcpPort"`
	Development         bool     `mapstructure:"development"                             env:"Development"`
	BasePath            string   `mapstructure:"basePath"            validate:"required" env:"BasePath"`
	DebugErrorsResponse bool     `mapstructure:"debugErrorsResponse"                     env:"DebugErrorsResponse"`
	IgnoreLogUrls       []string `mapstructure:"ignoreLogUrls"`
	Timeout             int      `mapstructure:"timeout"                                 env:"Timeout"`
	Host                string   `mapstructure:"host"                                    env:"Host"`
	Name                string   `mapstructure:"name"                                    env:"ShortTypeName"`
}

func (c *EchoHttpOptions) Address() string {
	return fmt.Sprintf("%s%s", c.Host, c.Port)
}

func (c *EchoHttpOptions) BasePathAddress() string {
	path, err := url.JoinPath(c.Address(), c.BasePath)
	if err != nil {
		return ""
	}
	return path
}

func ProvideConfig(environment environment.Environment) (*EchoHttpOptions, error) {
	return config.BindConfigKey[*EchoHttpOptions](optionName, environment)
}
