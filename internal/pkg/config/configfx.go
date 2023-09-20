package config

import (
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/config/environemnt"

	"go.uber.org/fx"
)

// Module provided to fxlog
// https://uber-go.github.io/fx/modules.html
var Module = fx.Module(
	"configfx",
	fx.Provide(func() environemnt.Environment {
		return environemnt.ConfigAppEnv()
	}),
)

var ModuleFunc = func(e environemnt.Environment) fx.Option {
	return fx.Module(
		"configfx",
		fx.Provide(func() environemnt.Environment {
			return environemnt.ConfigAppEnv(e)
		}),
	)
}
