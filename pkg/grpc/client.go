package grpc

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type grpcClient struct {
	conn *grpc.ClientConn
}

type GrpcClient interface {
	GetGrpcConnection() *grpc.ClientConn
	Close() error
}

func NewGrpcClient(config *GrpcConfig) (GrpcClient, error) {
	// Grpc Client to call Grpc Server
	//https://sahansera.dev/building-grpc-client-go/
	//https://github.com/open-telemetry/opentelemetry-go-contrib/blob/df16f32df86b40077c9c90d06f33c4cdb6dd5afa/instrumentation/google.golang.org/grpc/otelgrpc/example_interceptor_test.go
	conn, err := grpc.Dial(fmt.Sprintf("%s%s", config.Host, config.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	return &grpcClient{conn: conn}, nil
}

func (g *grpcClient) GetGrpcConnection() *grpc.ClientConn {
	return g.conn
}

func (g *grpcClient) Close() error {
	return g.conn.Close()
}
