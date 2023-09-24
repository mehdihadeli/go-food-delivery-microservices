package health

import (
	"go.uber.org/fx"
)

var Module = fx.Options( //nolint:gochecknoglobals
	fx.Provide(
		NewHealthService,
		NewHealthCheckEndpoint,
	),
	fx.Invoke(func(endpoint *HealthCheckEndpoint) {
		endpoint.RegisterEndpoints()
	}),
)
