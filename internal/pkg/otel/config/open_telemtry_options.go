package config

import (
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/config"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/config/environemnt"
	typeMapper "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/reflection/type_mappper"

	"github.com/iancoleman/strcase"
)

type OpenTelemetryOptions struct {
	Enabled               bool                   `mapstructure:"enabled"`
	ServiceName           string                 `mapstructure:"serviceName"`
	InstrumentationName   string                 `mapstructure:"instrumentationName"`
	Id                    int64                  `mapstructure:"id"`
	AlwaysOnSampler       bool                   `mapstructure:"alwaysOnSampler"`
	JaegerExporterOptions *JaegerExporterOptions `mapstructure:"jaegerExporterOptions"`
	ZipkinExporterOptions *ZipkinExporterOptions `mapstructure:"zipkinExporterOptions"`
	OTelMetricsOptions    *OTelMetricsOptions    `mapstructure:"otelMetricsOptions"`
	UseStdout             bool                   `mapstructure:"useStdout"`
}

type JaegerExporterOptions struct {
	OtlpEndpoint string `mapstructure:"otlpEndpoint"`
}

type ZipkinExporterOptions struct {
	Url string `mapstructure:"url"`
}

type OTelMetricsOptions struct {
	Host             string `mapstructure:"host"`
	Port             string `mapstructure:"port"`
	Name             string `mapstructure:"name"`
	MetricsRoutePath string `mapstructure:"metricsRoutePath"`
}

func ProvideOtelConfig(
	environment environemnt.Environment,
) (*OpenTelemetryOptions, error) {
	optionName := strcase.ToLowerCamel(
		typeMapper.GetTypeNameByT[OpenTelemetryOptions](),
	)

	return config.BindConfigKey[*OpenTelemetryOptions](optionName, environment)
}
