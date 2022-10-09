package grpc

import (
	"context"
	"emperror.dev/errors"
	"fmt"
	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcRecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpcCtxTags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpcError "github.com/mehdihadeli/store-golang-microservice-sample/pkg/grpc/interceptors/grpc_error"
	otelMetrics "github.com/mehdihadeli/store-golang-microservice-sample/pkg/grpc/interceptors/otel_metrics"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/metric"
	googleGrpc "google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
	"net"
	"time"
)

const (
	maxConnectionIdle = 5
	gRPCTimeout       = 15
	maxConnectionAge  = 5
	gRPCTime          = 10
)

type GrpcConfig struct {
	Port        string `mapstructure:"port" env:"Port"`
	Host        string `mapstructure:"host" env:"Host"`
	Development bool   `mapstructure:"development" env:"Development"`
}

type GrpcServer interface {
	RunGrpcServer(ctx context.Context, configGrpc ...func(grpcServer *googleGrpc.Server)) error
	GracefulShutdown()
	GetCurrentGrpcServer() *googleGrpc.Server
	GrpcServiceBuilder() *GrpcServiceBuilder
}

type grpcServer struct {
	server         *googleGrpc.Server
	config         *GrpcConfig
	log            logger.Logger
	serviceName    string
	serviceBuilder *GrpcServiceBuilder
}

func NewGrpcServer(config *GrpcConfig, logger logger.Logger, serviceName string, meter metric.Meter) *grpcServer {
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
		unaryServerInterceptors = append(unaryServerInterceptors, otelMetrics.UnaryServerInterceptor(meter, serviceName))
		streamServerInterceptors = append(streamServerInterceptors, otelMetrics.StreamServerInterceptor(meter, serviceName))
	}

	s := googleGrpc.NewServer(
		googleGrpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: maxConnectionIdle * time.Minute,
			Timeout:           gRPCTimeout * time.Second,
			MaxConnectionAge:  maxConnectionAge * time.Minute,
			Time:              gRPCTime * time.Minute,
		}),
		//https://github.com/open-telemetry/opentelemetry-go-contrib/tree/00b796d0cdc204fa5d864ec690b2ee9656bb5cfc/instrumentation/google.golang.org/grpc/otelgrpc
		//github.com/grpc-ecosystem/go-grpc-middleware
		googleGrpc.StreamInterceptor(grpcMiddleware.ChainStreamServer(
			streamServerInterceptors...,
		)),
		googleGrpc.UnaryInterceptor(grpcMiddleware.ChainUnaryServer(
			unaryServerInterceptors...,
		)),
	)

	return &grpcServer{server: s, config: config, log: logger, serviceName: serviceName, serviceBuilder: NewGrpcServiceBuilder(s)}
}

func (s *grpcServer) RunGrpcServer(ctx context.Context, configGrpc ...func(grpcServer *googleGrpc.Server)) error {
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

	go func() {
		for {
			select {
			case <-ctx.Done():
				s.log.Infof("%s is shutting down Grpc PORT: {%s}", s.serviceName, s.config.Port)
				s.GracefulShutdown()
				return
			}
		}
	}()

	s.log.Infof("[grpcServer.RunGrpcServer] Writer gRPC server is listening on port: %s", s.config.Port)

	err = s.server.Serve(l)

	if err != nil {
		s.log.Error(fmt.Sprintf("[grpcServer_RunGrpcServer.Serve] grpc server serve error: %+v", err))
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
