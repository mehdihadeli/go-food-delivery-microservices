package server

import (
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

const (
	maxConnectionIdle = 5
	gRPCTimeout       = 15
	maxConnectionAge  = 5
	gRPCTime          = 10
)

func (s Server) RunGrpcServer(registerServiceServer func(grpcServer *grpc.Server)) (error, func()) {

	l, err := net.Listen("tcp", s.Cfg.GRPC.Port)
	if err != nil {
		return errors.Wrap(err, "net.Listen"), nil
	}

	if registerServiceServer != nil {
		registerServiceServer(s.GrpcServer)
	}

	grpc_prometheus.Register(s.GrpcServer)

	if s.Cfg.GRPC.Development {
		reflection.Register(s.GrpcServer)
	}

	go func() {
		s.Log.Infof("Writer gRPC server is listening on port: %s", s.Cfg.GRPC.Port)
		s.Log.Fatal(s.GrpcServer.Serve(l))
	}()

	return nil, func() {
		l.Close()
		s.GrpcServer.GracefulStop()
	}
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
