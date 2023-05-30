package customEcho

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/labstack/gommon/log"
	"go.uber.org/fx"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/config"
)

// Module provided to fx
// https://uber-go.github.io/fx/modules.html
var Module = fx.Module(
	"customechofx",
	// - order is not important in provide
	// - provide can have parameter and will resolve if registered
	// - execute its func only if it requested
	fx.Provide(
		provideConfig,
		fx.Annotate(
			NewEchoHttpServer,
			fx.ParamTags(``, ``, `name:"meter" optional:"true"`),
		),
	),
	// - execute after registering all of our provided
	// - they execute by their orders
	// - invokes always execute its func compare to provides that only run when we request for them.
	// - return value will be discarded and can not be provided
	fx.Invoke(invokeHttp),
)

func provideConfig() (*EchoHttpConfig, error) {
	return config.BindConfigKey[*EchoHttpConfig]("http")
}

// we don't want to register any dependencies here, its func body should execute always even we don't request for that, so we should use `invoke`
func invokeHttp(lc fx.Lifecycle, echoServer EchoHttpServer) {
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
						Fatalf("(s.RunHttpServer) error in running server: {%v}", err)
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
				log.Errorf("error shutting down server: %v", err)
			} else {
				log.Info("server shutdown gracefully")
			}
			return nil
		},
	})
}
