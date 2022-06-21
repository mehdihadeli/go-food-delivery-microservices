package main

import (
	"context"
	"flag"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/products"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/shared/configurations"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/shared/server"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/shared/web/middlewares"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// @contact.name Mehdi Hadeli
// @contact.url https://github.com/mehdihadeli
func main() {
	flag.Parse()

	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatal(err)
	}

	appLogger := logger.NewAppLogger(cfg.Logger)
	appLogger.InitLogger()
	appLogger.WithName("catalogs-services")

	appLogger.Fatal(Run(appLogger, cfg))
}

func Run(log logger.Logger, cfg *config.Config) error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	s := server.NewServer(log, cfg)
	s.MiddlewareManager = middlewares.NewMiddlewareManager(log, cfg)

	ic := configurations.NewInfrastructureConfigurator(s)
	err, defers := ic.ConfigInfrastructures(ctx, cancel)

	if err != nil {
		return err
	}

	defer defers()

	pc := products.NewProductsModuleConfigurator(s)
	err = pc.ConfigureProductsModule()
	if err != nil {
		return err
	}

	go func() {
		if err := s.runHttpServer(); err != nil {
			s.Log.Errorf("(s.runHttpServer) err: {%v}", err)
			cancel()
		}
	}()
	s.Log.Infof("%s is listening on PORT: {%s}", GetMicroserviceName(s.Cfg), s.Cfg.Http.Port)

	<-ctx.Done()
	s.waitShootDown(waitShotDownDuration)

	if err := s.shutDownHealthCheckServer(ctx); err != nil {
		s.Log.Warnf("(shutDownHealthCheckServer) err: {%v}", err)
	}
	if err := s.Echo.Shutdown(ctx); err != nil {
		s.Log.Warnf("(Shutdown) err: {%v}", err)
	}

	<-s.doneCh
	s.Log.Infof("%s server exited properly", GetMicroserviceName(s.Cfg))

	return nil
}
