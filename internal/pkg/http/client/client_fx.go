package client

import "go.uber.org/fx"

// Module provided to fxlog
// https://uber-go.github.io/fx/modules.html
var Module = fx.Module("clientfx",
	// - order is not important in provide
	// - provide can have parameter and will resolve if registered
	// - execute its func only if it requested
	fx.Provide(NewHttpClient),
)
