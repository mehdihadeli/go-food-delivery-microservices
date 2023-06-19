package elasticsearch

import (
	"go.uber.org/fx"
)

// Module provided to fxlog
// https://uber-go.github.io/fx/modules.html
var Module = fx.Module("elasticfx",
	fx.Provide(provideConfig),
	fx.Provide(NewElasticClient),
)
