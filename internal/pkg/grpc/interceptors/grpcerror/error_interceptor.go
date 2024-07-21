package grpcError

import (
	"context"

	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/grpc/grpcerrors"

	"emperror.dev/errors"
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

		var grpcErr grpcerrors.GrpcErr

		// if error was not `grpcErr` we will convert the error to a `grpcErr`
		if ok := errors.As(err, &grpcErr); !ok {
			grpcErr = grpcerrors.ParseError(err)
		}

		if grpcErr != nil {
			return nil, grpcErr.ToGrpcResponseErr()
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

		var grpcErr grpcerrors.GrpcErr

		// if error was not `grpcErr` we will convert the error to a `grpcErr`
		if ok := errors.As(err, &grpcErr); !ok {
			grpcErr = grpcerrors.ParseError(err)
		}

		if grpcErr != nil {
			return grpcErr.ToGrpcResponseErr()
		}

		return err
	}
}
