package server

import (
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/pkg/errors"
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

func NewGrpcServer() *grpc.Server {
	grpcServer := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: maxConnectionIdle * time.Minute,
			Timeout:           gRPCTimeout * time.Second,
			MaxConnectionAge:  maxConnectionAge * time.Minute,
			Time:              gRPCTime * time.Minute,
		}),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_opentracing.UnaryServerInterceptor(),
			grpc_prometheus.UnaryServerInterceptor,
			grpc_recovery.UnaryServerInterceptor()),
		),
	)

	return grpcServer
}

func (s Server) RunGrpcServer(configGrpc func(grpcServer *grpc.Server)) error {

	l, err := net.Listen("tcp", s.Cfg.GRPC.Port)
	if err != nil {
		return errors.Wrap(err, "net.Listen")
	}

	if configGrpc != nil {
		configGrpc(s.GrpcServer)
	}

	grpc_prometheus.Register(s.GrpcServer)

	if s.Cfg.GRPC.Development {
		reflection.Register(s.GrpcServer)
	}

	s.Log.Infof("Writer gRPC server is listening on port: %s", s.Cfg.GRPC.Port)
	err = s.GrpcServer.Serve(l)
	s.Log.Fatal(err)

	return err
}

//func (s *Server) newCatalogsServiceGrpcServer() (error, func()) {
//	l, err := net.Listen("tcp", s.Cfg.GRPC.Port)
//	if err != nil {
//		return errors.Wrap(err, "net.Listen"), nil
//	}
//
//	grpcServer := grpc.NewServer(
//		grpc.KeepaliveParams(keepalive.ServerParameters{
//			MaxConnectionIdle: maxConnectionIdle * time.Minute,
//			Timeout:           gRPCTimeout * time.Second,
//			MaxConnectionAge:  maxConnectionAge * time.Minute,
//			Time:              gRPCTime * time.Minute,
//		}),
//		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
//			grpc_ctxtags.UnaryServerInterceptor(),
//			grpc_opentracing.UnaryServerInterceptor(),
//			grpc_prometheus.UnaryServerInterceptor,
//			grpc_recovery.UnaryServerInterceptor(),
//			s.Im.Logger,
//		),
//		),
//	)
//
//	productGrpcService := grpc_delivery.NewProductGrpcService(s.Log, s.Cfg, s.Validator, s.Mediator, s.Metrics)
//	product_service.RegisterProductsServiceServer(grpcServer, productGrpcService)
//	grpc_prometheus.Register(grpcServer)
//
//	if s.Cfg.GRPC.Development {
//		reflection.Register(grpcServer)
//	}
//
//	go func() {
//		s.Log.Infof("Writer gRPC server is listening on port: %s", s.cfg.GRPC.Port)
//		s.Log.Fatal(grpcServer.Serve(l))
//	}()
//
//	return nil, func() {
//		l.Close()
//		grpcServer.GracefulStop()
//	}
//}
