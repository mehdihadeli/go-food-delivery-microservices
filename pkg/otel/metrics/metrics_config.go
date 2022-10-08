package metrics

type OTelMetricsConfig struct {
	Host             string `mapstructure:"host"`
	Port             string `mapstructure:"port"`
	Name             string `mapstructure:"name"`
	MetricsRoutePath string `mapstructure:"metricsRoutePath"`
}
