package customEcho

import (
	"context"

	"go.uber.org/fx"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/custom_echo/config"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
)

// Module provided to fxlog
// https://uber-go.github.io/fx/modules.html
var Module = fx.Module(
	"customechofx",
	// - order is not important in provide
	// - provide can have parameter and will resolve if registered
	// - execute its func only if it requested
	fx.Provide(
		config.ProvideConfig,
		// https://uber-go.github.io/fx/value-groups/consume.html#with-annotated-functions
		// https://uber-go.github.io/fx/annotate.html
		fx.Annotate(
			NewEchoHttpServer,
			fx.ParamTags(``, ``, `optional:"true"`),
		),
	),
	// - execute after registering all of our provided
	// - they execute by their orders
	// - invokes always execute its func compare to provides that only run when we request for them.
	// - return value will be discarded and can not be provided
	fx.Invoke(registerHooks),
)

// we don't want to register any dependencies here, its func body should execute always even we don't request for that, so we should use `invoke`
func registerHooks(lc fx.Lifecycle, echoServer EchoHttpServer, logger logger.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			// https://github.com/uber-go/fx/blob/v1.20.0/app.go#L573
			// this ctx is just for startup dependencies setup and OnStart callbacks, and it has short timeout 15s, and it is not alive in whole lifetime app
			// if we need an app context which is alive until the app context done we should create it manually here

			go func() {
				if err := echoServer.RunHttpServer(); err != nil {
					// do a fatal for going to OnStop process
					logger.Fatalf(
						"(EchoHttpServer.RunHttpServer) error in running server: {%v}",
						err,
					)
				}
			}()
			echoServer.Logger().Infof(
				"%s is listening on Host:{%s} Http PORT: {%s}",
				echoServer.Cfg().Name,
				echoServer.Cfg().Host,
				echoServer.Cfg().Port,
			)

			return nil
		},
		OnStop: func(ctx context.Context) error {
			// https://github.com/uber-go/fx/blob/v1.20.0/app.go#L573
			// this ctx is just for stopping callbacks or OnStop callbacks, and it has short timeout 15s, and it is not alive in whole lifetime app
			if err := echoServer.GracefulShutdown(ctx); err != nil {
				echoServer.Logger().Errorf("error shutting down echo server: %v", err)
			} else {
				echoServer.Logger().Info("echo server shutdown gracefully")
			}
			return nil
		},
	})
}
