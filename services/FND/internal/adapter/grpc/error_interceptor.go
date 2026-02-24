package grpc

import (
	"context"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ErrorInterceptor creates a gRPC UnaryServerInterceptor that sanitizes internal errors.
// It prevents raw database errors or stack traces from leaking to the client.
func ErrorInterceptor(log *slog.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		res, err := handler(ctx, req)
		if err != nil {
			// Check if the error is already a formatted gRPC status error
			if _, ok := status.FromError(err); ok {
				// If it's a known Unauthenticated or InvalidArgument, let it pass
				st := status.Convert(err)
				if st.Code() != codes.Unknown && st.Code() != codes.Internal {
					return res, err
				}
			}

			// If it's a raw error (like GORM bubbling up a syntax/constraint error)
			// we log the actual error securely into our backend telemetry
			log.Error("Internal gRPC error during execution",
				slog.String("method", info.FullMethod),
				slog.String("error", err.Error()),
			)

			// And return a generic sanitized response to the frontend client
			return nil, status.Error(codes.Internal, "an internal server error occurred")
		}

		return res, nil
	}
}
