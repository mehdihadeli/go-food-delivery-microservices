package config

import "go.uber.org/fx"

// Module provided to fxlog
// https://uber-go.github.io/fx/modules.html
var Module = fx.Module(
	"configfx",
	fx.Provide(ConfigAppEnv),
)
