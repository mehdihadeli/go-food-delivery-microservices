package grpc

import (
	"google.golang.org/grpc"
)

type GrpcServiceBuilder struct {
	server *grpc.Server
}

func NewGrpcServiceBuilder(server *grpc.Server) *GrpcServiceBuilder {
	return &GrpcServiceBuilder{server: server}
}

func (r *GrpcServiceBuilder) RegisterRoutes(builder func(s *grpc.Server)) *GrpcServiceBuilder {
	builder(r.server)

	return r
}

func (r *GrpcServiceBuilder) Build() *grpc.Server {
	return r.server
}
