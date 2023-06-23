package config

import (
	"github.com/iancoleman/strcase"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/config"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/config/environemnt"
	typeMapper "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/reflection/type_mappper"
)

var optionName = strcase.ToLowerCamel(typeMapper.GetTypeNameByT[EventStoreDbOptions]())

type EventStoreDbOptions struct {
	ConnectionString string        `mapstructure:"connectionString"`
	Subscription     *Subscription `mapstructure:"subscription"`
}

type Subscription struct {
	Prefix         []string `mapstructure:"prefix"         validate:"required"`
	SubscriptionId string   `mapstructure:"subscriptionId" validate:"required"`
}

func ProvideConfig(environment environemnt.Environment) (*EventStoreDbOptions, error) {
	return config.BindConfigKey[*EventStoreDbOptions](optionName, environment)
}
