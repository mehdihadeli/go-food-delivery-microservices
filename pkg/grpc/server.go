package grpc

import (
	"emperror.dev/errors"
	"fmt"
	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcRecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpcCtxTags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpcOpentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	grpcPrometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"google.golang.org/grpc"
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
	Development bool   `mapstructure:"development" env:"Development"`
}

type GrpcServer interface {
	RunGrpcServer(configGrpc func(grpcServer *grpc.Server)) error
	GracefulShutdown()
	GetCurrentGrpcServer() *grpc.Server
}

type grpcServer struct {
	server *grpc.Server
	config *GrpcConfig
	log    logger.Logger
}

func NewGrpcServer(config *GrpcConfig, logger logger.Logger) *grpcServer {
	s := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: maxConnectionIdle * time.Minute,
			Timeout:           gRPCTimeout * time.Second,
			MaxConnectionAge:  maxConnectionAge * time.Minute,
			Time:              gRPCTime * time.Minute,
		}),
		grpc.UnaryInterceptor(grpcMiddleware.ChainUnaryServer(
			grpcCtxTags.UnaryServerInterceptor(),
			grpcOpentracing.UnaryServerInterceptor(),
			grpcPrometheus.UnaryServerInterceptor,
			grpcRecovery.UnaryServerInterceptor()),
		),
	)

	return &grpcServer{server: s, config: config, log: logger}
}

func (s *grpcServer) RunGrpcServer(configGrpc func(grpcServer *grpc.Server)) error {
	l, err := net.Listen("tcp", s.config.Port)
	if err != nil {
		return errors.WrapIf(err, "net.Listen")
	}

	if configGrpc != nil {
		configGrpc(s.server)
	}

	grpcPrometheus.Register(s.server)

	if s.config.Development {
		reflection.Register(s.server)
	}

	s.log.Infof("[grpcServer.RunGrpcServer] Writer gRPC server is listening on port: %s", s.config.Port)

	err = s.server.Serve(l)

	if err != nil {
		s.log.Error(fmt.Sprintf("[grpcServer_RunGrpcServer.Serve] grpc server serve error: %+v", err))
	}

	return err
}

func (s *grpcServer) GetCurrentGrpcServer() *grpc.Server {
	return s.server
}

func (s *grpcServer) GracefulShutdown() {
	s.server.Stop()
	s.server.GracefulStop()
}
