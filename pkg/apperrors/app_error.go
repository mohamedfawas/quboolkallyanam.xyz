package apperrors

import (
	"errors"

	"google.golang.org/grpc/codes"
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
