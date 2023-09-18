package health

import (
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		NewHealthService,
		NewHealthCheckEndpoint,
	),
	fx.Invoke(func(endpoint *HealthCheckEndpoint) {
		endpoint.RegisterEndpoints()
	}),
)
