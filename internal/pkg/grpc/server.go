package grpc

import (
	"fmt"
	"net"
	"time"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/grpc/config"
	grpcError "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/grpc/interceptors/grpc_error"
	otelMetrics "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/grpc/interceptors/otel_metrics"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"

	"emperror.dev/errors"
	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcRecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpcCtxTags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/metric"
	googleGrpc "google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

const (
	maxConnectionIdle = 5
	gRPCTimeout       = 15
	maxConnectionAge  = 5
	gRPCTime          = 10
)

type GrpcServer interface {
	RunGrpcServer(configGrpc ...func(grpcServer *googleGrpc.Server)) error
	GracefulShutdown()
	GetCurrentGrpcServer() *googleGrpc.Server
	GrpcServiceBuilder() *GrpcServiceBuilder
}

type grpcServer struct {
	server         *googleGrpc.Server
	config         *config.GrpcOptions
	log            logger.Logger
	serviceName    string
	serviceBuilder *GrpcServiceBuilder
}

func NewGrpcServer(
	config *config.GrpcOptions,
	logger logger.Logger,
	meter metric.Meter,
) GrpcServer {
	unaryServerInterceptors := []googleGrpc.UnaryServerInterceptor{
		otelgrpc.UnaryServerInterceptor(),
		grpcError.UnaryServerInterceptor(),
		grpcCtxTags.UnaryServerInterceptor(),
		grpcRecovery.UnaryServerInterceptor(),
	}
	streamServerInterceptors := []googleGrpc.StreamServerInterceptor{
		otelgrpc.StreamServerInterceptor(),
		grpcError.StreamServerInterceptor(),
	}

	if meter != nil {
		unaryServerInterceptors = append(
			unaryServerInterceptors,
			otelMetrics.UnaryServerInterceptor(meter, config.Name),
		)
		streamServerInterceptors = append(
			streamServerInterceptors,
			otelMetrics.StreamServerInterceptor(meter, config.Name),
		)
	}

	s := googleGrpc.NewServer(
		googleGrpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: maxConnectionIdle * time.Minute,
			Timeout:           gRPCTimeout * time.Second,
			MaxConnectionAge:  maxConnectionAge * time.Minute,
			Time:              gRPCTime * time.Minute,
		}),
		// https://github.com/open-telemetry/opentelemetry-go-contrib/tree/00b796d0cdc204fa5d864ec690b2ee9656bb5cfc/instrumentation/google.golang.org/grpc/otelgrpc
		// github.com/grpc-ecosystem/go-grpc-middleware
		googleGrpc.StreamInterceptor(grpcMiddleware.ChainStreamServer(
			streamServerInterceptors...,
		)),
		googleGrpc.UnaryInterceptor(grpcMiddleware.ChainUnaryServer(
			unaryServerInterceptors...,
		)),
	)

	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(s, healthServer)
	healthServer.SetServingStatus(config.Name, grpc_health_v1.HealthCheckResponse_SERVING)

	return &grpcServer{
		server:         s,
		config:         config,
		log:            logger,
		serviceName:    config.Name,
		serviceBuilder: NewGrpcServiceBuilder(s),
	}
}

func (s *grpcServer) RunGrpcServer(
	configGrpc ...func(grpcServer *googleGrpc.Server),
) error {
	l, err := net.Listen("tcp", s.config.Port)
	if err != nil {
		return errors.WrapIf(err, "net.Listen")
	}

	if len(configGrpc) > 0 {
		grpcFunc := configGrpc[0]
		if grpcFunc != nil {
			grpcFunc(s.server)
		}
	}

	if s.config.Development {
		reflection.Register(s.server)
	}

	s.log.Infof(
		"[grpcServer.RunGrpcServer] Writer gRPC server is listening on port: %s",
		s.config.Port,
	)

	err = s.server.Serve(l)

	if err != nil {
		s.log.Error(
			fmt.Sprintf("[grpcServer_RunGrpcServer.Serve] grpc server serve error: %+v", err),
		)
	}

	return err
}

func (s *grpcServer) GrpcServiceBuilder() *GrpcServiceBuilder {
	return s.serviceBuilder
}

func (s *grpcServer) GetCurrentGrpcServer() *googleGrpc.Server {
	return s.server
}

func (s *grpcServer) GracefulShutdown() {
	s.server.Stop()
	s.server.GracefulStop()
}
