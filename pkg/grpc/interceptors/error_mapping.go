package interceptors

import (
	"context"
	"errors"

	appErrors "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/errors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func UnaryErrorInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (resp interface{}, err error) {
		resp, err = handler(ctx, req)
		if err == nil {
			return resp, nil
		}

		var st *status.Status
		switch {
		case errors.Is(err, appErrors.ErrInvalidCredentials):
			st = status.New(codes.Unauthenticated, "Invalid credentials")

		// TODO: Add more error cases here
		default:
			if s, ok := status.FromError(err); ok {
				st = s
			} else {
				st = status.New(codes.Internal, "Internal server error")
			}
		}

		// Log full error server-side
		if st.Code() == codes.Internal {
			logger.Log.Error("gRPC handler error",
				zap.String("method", info.FullMethod),
				zap.Error(err),
				zap.String("mapped_message", st.Message()),
				zap.String("stack", zap.Stack("stack").String),
			)
		} else {
			logger.Log.Infow("gRPC client error",
				"method", info.FullMethod,
				"error", err.Error(),
				"mapped_code", st.Code(),
				"mapped_message", st.Message(),
			)
		}

		return nil, st.Err()
	}
}
