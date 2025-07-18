package interceptors

import (
	"context"
	"errors"

	appErrors "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/errors"
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

		// Map application errors to gRPC status codes
		st := mapAppErrorToGRPCStatus(err)

		return nil, st.Err()
	}
}
func mapAppErrorToGRPCStatus(err error) *status.Status {
	switch {
	// Authentication errors
	case errors.Is(err, appErrors.ErrInvalidCredentials):
		return status.New(codes.Unauthenticated, "Invalid credentials")
	case errors.Is(err, appErrors.ErrInvalidToken):
		return status.New(codes.Unauthenticated, "Invalid token")
	case errors.Is(err, appErrors.ErrExpiredToken):
		return status.New(codes.Unauthenticated, "Token expired")
	case errors.Is(err, appErrors.ErrInvalidRefreshToken):
		return status.New(codes.Unauthenticated, "Invalid refresh token")
	case errors.Is(err, appErrors.ErrUnauthorized):
		return status.New(codes.Unauthenticated, "Unauthorized")

	// Permission errors
	case errors.Is(err, appErrors.ErrForbidden):
		return status.New(codes.PermissionDenied, "Access denied")
	case errors.Is(err, appErrors.ErrAccountDisabled):
		return status.New(codes.PermissionDenied, "Account disabled")
	case errors.Is(err, appErrors.ErrAccountBlocked):
		return status.New(codes.PermissionDenied, "Account blocked")
	case errors.Is(err, appErrors.ErrAdminAccountDisabled):
		return status.New(codes.PermissionDenied, "Admin account disabled")

	// Not found errors
	case errors.Is(err, appErrors.ErrUserNotFound):
		return status.New(codes.NotFound, "User not found")
	case errors.Is(err, appErrors.ErrAdminNotFound):
		return status.New(codes.NotFound, "Admin not found")
	case errors.Is(err, appErrors.ErrPendingRegistrationNotFound):
		return status.New(codes.NotFound, "Registration not found")
	case errors.Is(err, appErrors.ErrOTPNotFound):
		return status.New(codes.NotFound, "OTP not found")
	case errors.Is(err, appErrors.ErrPaymentNotFound):
		return status.New(codes.NotFound, "Payment not found")
	case errors.Is(err, appErrors.ErrSubscriptionPlanNotFound):
		return status.New(codes.NotFound, "Subscription plan not found")

	// Conflict errors
	case errors.Is(err, appErrors.ErrEmailAlreadyExists):
		return status.New(codes.AlreadyExists, "Email already exists")
	case errors.Is(err, appErrors.ErrPhoneAlreadyExists):
		return status.New(codes.AlreadyExists, "Phone already exists")

	// Validation errors
	case errors.Is(err, appErrors.ErrInvalidEmail):
		return status.New(codes.InvalidArgument, "Invalid email")
	case errors.Is(err, appErrors.ErrInvalidPhoneNumber):
		return status.New(codes.InvalidArgument, "Invalid phone number")
	case errors.Is(err, appErrors.ErrInvalidPassword):
		return status.New(codes.InvalidArgument, "Invalid password")
	case errors.Is(err, appErrors.ErrInvalidInput):
		return status.New(codes.InvalidArgument, "Invalid input")
	case errors.Is(err, appErrors.ErrInvalidOTP):
		return status.New(codes.InvalidArgument, "Invalid OTP")
	case errors.Is(err, appErrors.ErrPaymentSignatureInvalid):
		return status.New(codes.InvalidArgument, "Payment verification failed")

	// Payment state errors
	case errors.Is(err, appErrors.ErrPaymentAlreadyCompleted):
		return status.New(codes.FailedPrecondition, "Payment already completed")
	case errors.Is(err, appErrors.ErrPaymentExpired):
		return status.New(codes.FailedPrecondition, "Payment has expired")
	case errors.Is(err, appErrors.ErrSubscriptionPlanNotActive):
		return status.New(codes.FailedPrecondition, "Subscription plan is not active")

	// Account state errors
	case errors.Is(err, appErrors.ErrAccountNotVerified):
		return status.New(codes.FailedPrecondition, "Account not verified")

	// Payment processing errors
	case errors.Is(err, appErrors.ErrRazorpayOrderCreation):
		return status.New(codes.Unavailable, "Payment service temporarily unavailable")
	case errors.Is(err, appErrors.ErrPaymentProcessingFailed):
		return status.New(codes.Internal, "Payment processing failed")

	default:
		// Check if it's already a gRPC status error
		if st, ok := status.FromError(err); ok {
			return st
		}
		// Default to internal error
		return status.New(codes.Internal, "Internal server error")
	}
}
