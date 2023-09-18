package config

import (
	"fmt"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/config"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/config/environemnt"
	typeMapper "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/reflection/type_mappper"

	"github.com/iancoleman/strcase"
)

// https://developers.eventstore.com/clients/dotnet/21.2/#connect-to-eventstoredb
// https://developers.eventstore.com/clients/http-api/v5
// https://developers.eventstore.com/clients/grpc/
// https://developers.eventstore.com/server/v20.10/networking.html#http-configuration

type EventStoreDbOptions struct {
	Host    string `mapstructure:"host"`
	TcpPort int    `mapstructure:"tcpPort"`
	// HTTP is the primary protocol for EventStoreDB. It is used in gRPC communication and HTTP APIs (management, gossip and diagnostics).
	HttpPort     int           `mapstructure:"httpPort"`
	Subscription *Subscription `mapstructure:"subscription"`
}

// https://developers.eventstore.com/server/v20.10/networking.html#http-configuration
// https://developers.eventstore.com/clients/grpc/#connection-string

func (e *EventStoreDbOptions) GrpcEndPoint() string {
	return fmt.Sprintf("esdb://%s:%d?tls=false", e.Host, e.HttpPort)
}

// https://developers.eventstore.com/clients/dotnet/21.2/#connect-to-eventstoredb
// https://developers.eventstore.com/server/v20.10/networking.html#external

func (e *EventStoreDbOptions) TcpEndPoint() string {
	return fmt.Sprintf("tcp://%s:%d?tls=false", e.Host, e.TcpPort)
}

// https://developers.eventstore.com/server/v20.10/networking.html#http-configuration
// https://developers.eventstore.com/clients/http-api/v5

func (e *EventStoreDbOptions) HttpEndPoint() string {
	return fmt.Sprintf("http://%s:%d", e.Host, e.Host)
}

type Subscription struct {
	Prefix         []string `mapstructure:"prefix"         validate:"required"`
	SubscriptionId string   `mapstructure:"subscriptionId" validate:"required"`
}

func ProvideConfig(environment environemnt.Environment) (*EventStoreDbOptions, error) {
	optionName := strcase.ToLowerCamel(typeMapper.GetTypeNameByT[EventStoreDbOptions]())
	return config.BindConfigKey[*EventStoreDbOptions](optionName, environment)
}
