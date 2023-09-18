package queue

import (
	"go.uber.org/fx"
)

var (
	ClientModule    = fx.Module("queue-client", ClientProviders, ClientInvokes)
	ClientProviders = fx.Options(
		fx.Provide(NewServeMux),
		fx.Provide(NewClient),
	)
	ClientInvokes = fx.Options(
		fx.Invoke(HookClient),
	)

	WorkerModule    = fx.Module("queue-worker", ClientModule, WorkerProviders, WorkerInvokes)
	WorkerProviders = fx.Options(
		fx.Provide(NewServer),
	)
	WorkerInvokes = fx.Options(
		fx.Invoke(HookServer),
	)
)
