package metrics

import (
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/config"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/config/environemnt"
	typeMapper "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/reflection/type_mappper"

	"github.com/iancoleman/strcase"
)

type OTLPProvider struct {
	Name         string            `mapstructure:"name"`
	Enabled      bool              `mapstructure:"enabled"`
	OTLPEndpoint string            `mapstructure:"otlpEndpoint"`
	OTLPHeaders  map[string]string `mapstructure:"otlpHeaders"`
}

type MetricsOptions struct {
	Host                      string         `mapstructure:"host"`
	Port                      string         `mapstructure:"port"`
	ServiceName               string         `mapstructure:"serviceName"`
	Version                   string         `mapstructure:"version"`
	MetricsRoutePath          string         `mapstructure:"metricsRoutePath"`
	EnableHostMetrics         bool           `mapstructure:"enableHostMetrics"`
	UseStdout                 bool           `mapstructure:"useStdout"`
	InstrumentationName       string         `mapstructure:"instrumentationName"`
	UseOTLP                   bool           `mapstructure:"useOTLP"`
	OTLPProviders             []OTLPProvider `mapstructure:"otlpProviders"`
	ElasticApmExporterOptions *OTLPProvider  `mapstructure:"elasticApmExporterOptions"`
	UptraceExporterOptions    *OTLPProvider  `mapstructure:"uptraceExporterOptions"`
}

func ProvideMetricsConfig(
	environment environemnt.Environment,
) (*MetricsOptions, error) {
	optionName := strcase.ToLowerCamel(
		typeMapper.GetTypeNameByT[MetricsOptions](),
	)

	return config.BindConfigKey[*MetricsOptions](optionName, environment)
}
