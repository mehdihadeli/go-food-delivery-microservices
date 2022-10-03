package grpcError

import (
	"context"
	"fmt"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/grpc/grpcErrors"
	"google.golang.org/grpc"
)

// UnaryServerInterceptor returns a problem-detail error to client
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		resp, err := handler(ctx, req)
		if err != nil {
			grpcErr := grpcErrors.ParseError(err)

			if grpcErr != nil {
				return nil, grpcErr.ToGrpcResponseErr()

			} else {
				prb := grpcErrors.NewInternalServerGrpcError(err.Error(), fmt.Sprintf("%+v\n", err))
				return nil, prb.ToGrpcResponseErr()
			}
		}

		return resp, err
	}
}

// StreamServerInterceptor returns a problem-detail error to client.
func StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		err := handler(srv, ss)

		if err != nil {
			grpcErr := grpcErrors.ParseError(err)

			if grpcErr != nil {
				return grpcErr.ToGrpcResponseErr()

			} else {
				prb := grpcErrors.NewInternalServerGrpcError(err.Error(), fmt.Sprintf("%+v\n", err))
				return prb.ToGrpcResponseErr()
			}
		}
		return err
	}
}
