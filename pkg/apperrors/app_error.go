package apperrors

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AppError struct {
	Err            error  // Underlying error
	Code           string // Machine readable error code
	HTTPStatusCode int
	GRPCStatusCode codes.Code
	PublicMsg      string // User Friendly msg
}

func (e *AppError) Error() string { return e.Err.Error() }
func (e *AppError) Unwrap() error { return e.Err }

func IsAppError(err error) bool {
	var ae *AppError

	return errors.As(err, &ae)
}

func ShouldLogError(err error) bool { // Used in gateway service only, never use in internal grpc services.
	if err == nil {
		return false
	}

	if IsAppError(err) {
		return false
	}

	if _, ok := status.FromError(err); ok {
		return false
	}

	return true
}


