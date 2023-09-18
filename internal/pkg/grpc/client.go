package grpc

import (
	"fmt"
	"time"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/grpc/config"

	"emperror.dev/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
)

type grpcClient struct {
	conn *grpc.ClientConn
}

type GrpcClient interface {
	GetGrpcConnection() *grpc.ClientConn
	Close() error
	// WaitForAvailableConnection waiting for grpc endpoint becomes ready in the given timeout
	WaitForAvailableConnection() error
}

func NewGrpcClient(config *config.GrpcOptions) (GrpcClient, error) {
	// Grpc Client to call Grpc Server
	// https://sahansera.dev/building-grpc-client-go/
	// https://github.com/open-telemetry/opentelemetry-go-contrib/blob/df16f32df86b40077c9c90d06f33c4cdb6dd5afa/instrumentation/google.golang.org/grpc/otelgrpc/example_interceptor_test.go
	conn, err := grpc.Dial(fmt.Sprintf("%s%s", config.Host, config.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	return &grpcClient{conn: conn}, err
}

func (g *grpcClient) GetGrpcConnection() *grpc.ClientConn {
	return g.conn
}

func (g *grpcClient) Close() error {
	return g.conn.Close()
}

func (g *grpcClient) WaitForAvailableConnection() error {
	timeout := time.Second * 20

	err := waitUntilConditionMet(func() bool {
		return g.conn.GetState() == connectivity.Ready
	}, timeout)

	state := g.conn.GetState()
	fmt.Println(fmt.Sprintf("grpc state is:%s", state))
	return err
}

func waitUntilConditionMet(conditionToMet func() bool, timeout ...time.Duration) error {
	timeOutTime := 20 * time.Second
	if len(timeout) >= 0 && timeout != nil {
		timeOutTime = timeout[0]
	}

	startTime := time.Now()
	timeOutExpired := false
	meet := conditionToMet()
	for meet == false {
		if timeOutExpired {
			return errors.New("grpc connection could not be established in the given timeout.")
		}
		time.Sleep(time.Second * 2)
		meet = conditionToMet()
		timeOutExpired = time.Now().Sub(startTime) > timeOutTime
	}

	return nil
}
