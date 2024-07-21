package elasticsearch

import (
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/config"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/config/environment"
	typeMapper "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/reflection/typemapper"

	"github.com/iancoleman/strcase"
)

var optionName = strcase.ToLowerCamel(typeMapper.GetGenericTypeNameByT[ElasticOptions]())

type ElasticOptions struct {
	URL string `mapstructure:"url"`
}

func provideConfig(environment environment.Environment) (*ElasticOptions, error) {
	return config.BindConfigKey[*ElasticOptions](optionName, environment)
}
