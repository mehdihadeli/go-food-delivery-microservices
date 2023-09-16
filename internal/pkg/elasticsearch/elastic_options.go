package elasticsearch

import (
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/config"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/config/environemnt"
	typeMapper "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/reflection/type_mappper"

	"github.com/iancoleman/strcase"
)

var optionName = strcase.ToLowerCamel(typeMapper.GetTypeNameByT[ElasticOptions]())

type ElasticOptions struct {
	URL string `mapstructure:"url"`
}

func provideConfig(environment environemnt.Environment) (*ElasticOptions, error) {
	return config.BindConfigKey[*ElasticOptions](optionName, environment)
}
