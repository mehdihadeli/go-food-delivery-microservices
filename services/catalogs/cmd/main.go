package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/products"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/shared/configurations"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/shared/constants"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/shared/server"
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

	ic := configurations.NewInfrastructureConfigurator(s)
	err, infrastructure, defers := ic.ConfigInfrastructures(ctx, cancel)

	if err != nil {
		return err
	}

	defer defers()

	pc := products.NewProductsModuleConfigurator(infrastructure)
	err = pc.ConfigureProductsModule()
	if err != nil {
		return err
	}

	s.RunMetrics(cancel)

	healthCleanup := s.RunHealthCheck(ctx)
	defer healthCleanup()

	go func() {
		deliveryType := s.Cfg.DeliveryType
		switch deliveryType {
		case "http":

			if err := s.RunHttpServer(); err != nil {
				s.Log.Errorf("(s.runHttpServer) err: {%v}", err)
				cancel()
			}
			s.Log.Infof("%s is listening on Http PORT: {%s}", configurations.GetMicroserviceName(s.Cfg), s.Cfg.Http.Port)

		case "grpc":
			err, grpcderfer := s.RunGrpcServer(nil)
			if err != nil {
				return
			}
			s.Log.Infof("%s is listening on Grpc PORT: {%s}", configurations.GetMicroserviceName(s.Cfg), s.Cfg.GRPC.Port)
			defer grpcderfer()
		default:
			fmt.Sprintf("server type %s is not supported", deliveryType)
			//panic()
		}

		if err := s.RunHttpServer(); err != nil {
			s.Log.Errorf("(s.runHttpServer) err: {%v}", err)
			cancel()
		}
	}()

	<-ctx.Done()
	s.WaitShootDown(constants.WaitShotDownDuration)

	if err := s.Echo.Shutdown(ctx); err != nil {
		s.Log.Warnf("(Shutdown) err: {%v}", err)
	}

	<-s.DoneCh
	s.Log.Infof("%s server exited properly", configurations.GetMicroserviceName(s.Cfg))

	return nil
}
