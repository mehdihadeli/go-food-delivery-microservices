package customEcho

import (
	"fmt"
	"net/url"
)

type EchoHttpConfig struct {
	Port                string   `mapstructure:"port" validate:"required" env:"Port"`
	Development         bool     `mapstructure:"development" env:"Development"`
	BasePath            string   `mapstructure:"basePath" validate:"required" env:"BasePath"`
	DebugErrorsResponse bool     `mapstructure:"debugErrorsResponse" env:"DebugErrorsResponse"`
	IgnoreLogUrls       []string `mapstructure:"ignoreLogUrls"`
	Timeout             int      `mapstructure:"timeout" env:"Timeout"`
	Host                string   `mapstructure:"host" env:"Host"`
}

func (c *EchoHttpConfig) Address() string {
	return fmt.Sprintf("%s%s", c.Host, c.Port)
}

func (c *EchoHttpConfig) BasePathAddress() string {
	path, err := url.JoinPath(c.Address(), c.BasePath)
	if err != nil {
		return ""
	}
	return path
}
