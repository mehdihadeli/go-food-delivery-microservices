package config

import (
	"go.uber.org/fx"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/constants"
)

// Module provided to fxlog
// https://uber-go.github.io/fx/modules.html
var Module = fx.Module(
	"configfx",
	fx.Provide(ConfigAppEnv),
)

var TestModule = fx.Module(
	"configfx",
	fx.Supply(ConfigAppEnv(constants.Test)),
)
