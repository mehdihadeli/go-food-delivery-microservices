package server

import (
	"context"
	"fmt"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/constants"
	grpcServer "github.com/mehdihadeli/store-golang-microservice-sample/pkg/grpc"
	customEcho "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/custom_echo"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/shared/configurations/orders"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/shared/web"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	Log          logger.Logger
	Cfg          *config.Config
	GrpcServer   grpcServer.GrpcServer
	EchoServer   customEcho.EchoHttpServer
	HealthServer *http.Server
	doneCh       chan struct{}
}

func NewServer(log logger.Logger, cfg *config.Config) *Server {
	g := grpcServer.NewGrpcServer(cfg.GRPC, log)
	h := customEcho.NewEchoHttpServer(cfg.Http, log)
	return &Server{Log: log, Cfg: cfg, EchoServer: h, GrpcServer: g, HealthServer: NewHealthCheckServer(cfg)}
}

func (s *Server) Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	ordersConfigurator := orders.NewOrdersServiceConfigurator(s.Log, s.Cfg, s.EchoServer, s.GrpcServer)
	err, ordersCleanup := ordersConfigurator.ConfigureCatalogsService(ctx)
	if err != nil {
		return err
	}
	defer ordersCleanup()

	deliveryType := s.Cfg.DeliveryType

	s.RunMetrics(cancel)

	healthCleanup := s.RunHealthCheck(ctx)
	defer healthCleanup()

	switch deliveryType {
	case "http":
		go func() {
			if err := s.EchoServer.RunHttpServer(nil); err != nil {
				s.Log.Errorf("(s.RunHttpServer) err: {%v}", err)
				cancel()
			}
		}()
		s.Log.Infof("%s is listening on Http PORT: {%s}", web.GetMicroserviceName(s.Cfg), s.Cfg.Http.Port)

	case "grpc":
		go func() {
			if err := s.GrpcServer.RunGrpcServer(nil); err != nil {
				s.Log.Errorf("(s.RunGrpcServer) err: {%v}", err)
				cancel()
			}
		}()
		s.Log.Infof("%s is listening on Grpc PORT: {%s}", web.GetMicroserviceName(s.Cfg), s.Cfg.GRPC.Port)
	default:
		panic(fmt.Sprintf("server type %s is not supported", deliveryType))
	}

	//mongoProjection := NewOrderProjection(s.log, db, mongoRepository, s.cfg)
	//elasticProjection := elastic_projection.NewElasticProjection(s.log, db, elasticRepository, s.cfg)
	//
	//go func() {
	//	err := mongoProjection.Subscribe(ctx, []string{s.cfg.Subscriptions.OrderPrefix}, s.cfg.Subscriptions.PoolSize, mongoProjection.ProcessEvents)
	//	if err != nil {
	//		s.log.Errorf("(orderProjection.Subscribe) err: {%v}", err)
	//		cancel()
	//	}
	//}()
	//
	//go func() {
	//	err := elasticProjection.Subscribe(ctx, []string{s.cfg.Subscriptions.OrderPrefix}, s.cfg.Subscriptions.PoolSize, elasticProjection.ProcessEvents)
	//	if err != nil {
	//		s.log.Errorf("(elasticProjection.Subscribe) err: {%v}", err)
	//		cancel()
	//	}
	//}()

	<-ctx.Done()
	s.waitForShootDown(constants.WaitShotDownDuration)

	switch deliveryType {
	case "http":
		s.Log.Infof("%s is shutting down Http PORT: {%s}", web.GetMicroserviceName(s.Cfg), s.Cfg.Http.Port)
		if err := s.EchoServer.GracefulShutdown(ctx); err != nil {
			s.Log.Warnf("(Shutdown) err: {%v}", err)
		}
	case "grpc":
		s.Log.Infof("%s is shutting down Grpc PORT: {%s}", web.GetMicroserviceName(s.Cfg), s.Cfg.GRPC.Port)
		s.GrpcServer.GracefulShutdown()
	}

	<-s.doneCh
	s.Log.Infof("%s server exited properly", web.GetMicroserviceName(s.Cfg))

	return nil
}

func (s *Server) waitForShootDown(duration time.Duration) {
	go func() {
		time.Sleep(duration)
		s.doneCh <- struct{}{}
	}()
}
