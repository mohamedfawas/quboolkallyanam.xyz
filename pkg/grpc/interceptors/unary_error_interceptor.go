package interceptors

import (
	"context"
	"errors"
	"strconv"

	appErrors "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apperrors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func UnaryErrorInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		resp, err := handler(ctx, req)
		if err == nil {
			return resp, nil
		}

		var ae *appErrors.AppError
		if errors.As(err, &ae) {
			st, detErr := status.New(ae.GRPCStatusCode, ae.PublicMsg).
				WithDetails(&errdetails.ErrorInfo{
					Reason: ae.Code,
					Domain: constants.ServiceChat,
					Metadata: map[string]string{
						constants.HTTPStatusCode: strconv.Itoa(ae.HTTPStatusCode), 
						constants.UserFriendlyMessage:     ae.PublicMsg,                    
					},
				})
			if detErr != nil {
				// fallback to controlled public message
				return nil, status.Error(ae.GRPCStatusCode, ae.PublicMsg)
			}
			return nil, st.Err()
		}

		return nil, status.Error(codes.Internal, constants.InteralServerErrorMessage)
	}
}
