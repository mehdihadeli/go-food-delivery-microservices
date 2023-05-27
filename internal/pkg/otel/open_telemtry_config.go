package otel

type OpenTelemetryConfig struct {
	Enabled              bool                  `mapstructure:"enabled"`
	ServiceName          string                `mapstructure:"serviceName"`
	InstrumetationName   string                `mapstructure:"instrumetationName"`
	Id                   int64                 `mapstructure:"id"`
	AlwaysOnSampler      bool                  `mapstructure:"alwaysOnSampler"`
	JaegerExporterConfig *JaegerExporterConfig `mapstructure:"jaegerExporterConfig"`
	ZipkinExporterConfig *ZipkinExporterConfig `mapstructure:"zipkinExporterConfig"`
	UseStdout            bool                  `mapstructure:"useStdout"`
}

type JaegerExporterConfig struct {
	AgentHost string `mapstructure:"agentHost"`
	AgentPort string `mapstructure:"agentPort"`
}

type ZipkinExporterConfig struct {
	Url string `mapstructure:"url"`
}
