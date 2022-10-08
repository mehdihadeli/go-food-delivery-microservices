package customEcho

type EchoHttpConfig struct {
	Port                string   `mapstructure:"port" validate:"required" env:"Port"`
	Development         bool     `mapstructure:"development" env:"Development"`
	BasePath            string   `mapstructure:"basePath" validate:"required" env:"BasePath"`
	DebugErrorsResponse bool     `mapstructure:"debugErrorsResponse" env:"DebugErrorsResponse"`
	IgnoreLogUrls       []string `mapstructure:"ignoreLogUrls"`
	Timeout             int      `mapstructure:"timeout" env:"Timeout"`
	Host                string   `mapstructure:"host" env:"Host"`
}
