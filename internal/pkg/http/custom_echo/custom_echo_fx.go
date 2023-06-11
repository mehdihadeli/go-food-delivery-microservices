package customEcho

import (
	"context"
	"errors"
	"net/http"
	"time"

	"go.uber.org/fx"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/custom_echo/config"
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
		//// https://uber-go.github.io/fx/value-groups/consume.html#with-annotated-functions
		//// https://uber-go.github.io/fx/annotate.html
		//fxlog.Annotate(
		//	NewEchoHttpServer,
		//	fxlog.ParamTags(``, ``, `name:"meter" optional:"true"`),
		//),
		// https://uber-go.github.io/fx/parameter-objects.html#using-parameter-objects
		NewEchoHttpServer,
	),
	// - execute after registering all of our provided
	// - they execute by their orders
	// - invokes always execute its func compare to provides that only run when we request for them.
	// - return value will be discarded and can not be provided
	fx.Invoke(registerHooks),
)

// we don't want to register any dependencies here, its func body should execute always even we don't request for that, so we should use `invoke`
func registerHooks(lc fx.Lifecycle, echoServer EchoHttpServer) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			// Start server in a separate goroutine, this way when the server is shutdown "s.e.Start" will
			// return promptly, and the call to "s.e.Shutdown" is the one that will wait for all other
			// resources to be properly freed. If it was the other way around, the application would just
			// exit without gracefully shutting down the server.
			// For more details: https://medium.com/@momchil.dev/proper-http-shutdown-in-go-bd3bfaade0f2
			go func() {
				if err := echoServer.RunHttpServer(ctx, nil); !errors.Is(
					err,
					http.ErrServerClosed,
				) {
					echoServer.Logger().
						Fatalf("(EchoHttpServer.RunHttpServer) error in running server: {%v}", err)
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
			ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
			defer cancel()
			if err := echoServer.GracefulShutdown(ctx); err != nil {
				echoServer.Logger().Errorf("error shutting down server: %v", err)
			} else {
				echoServer.Logger().Info("server shutdown gracefully")
			}
			return nil
		},
	})
}
