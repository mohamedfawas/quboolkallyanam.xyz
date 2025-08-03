package apperrors

import (
	"errors"
	"net/http"

	"google.golang.org/grpc/codes"
)

// ///////////////// Common Errors ////////////////////////////////////////////////////////////
var (
	ErrUserContextMissing = &AppError{
		Err:            errors.New("user context missing"),
		Code:           "UNAUTHORIZED",
		HTTPStatusCode: http.StatusInternalServerError,
		GRPCStatusCode: codes.Internal,
		PublicMsg:      "Something went wrong. Please try again later.",
	}
	ErrRequestIDContextMissing = &AppError{
		Err:            errors.New("request ID context missing"),
		Code:           "INTERNAL_SERVER_ERROR",
		HTTPStatusCode: http.StatusInternalServerError,
		GRPCStatusCode: codes.Internal,
		PublicMsg:      "Something went wrong. Please try again later.",
	}
	ErrBindingJSON = &AppError{
		Err:            errors.New("binding JSON failed"),
		Code:           "BAD_REQUEST",
		HTTPStatusCode: http.StatusBadRequest,
		GRPCStatusCode: codes.InvalidArgument,
		PublicMsg:      "Invalid request body.",
	}
	ErrUserNotFound = &AppError{
		Err:            errors.New("user not found"),
		Code:           "USER_NOT_FOUND",
		HTTPStatusCode: http.StatusNotFound,
		GRPCStatusCode: codes.NotFound,
		PublicMsg:      "Please update your Profile and try again."}
	ErrPartnerUserNotFound = &AppError{
		Err:            errors.New("partner user not found"),
		Code:           "PARTNER_USER_NOT_FOUND",
		HTTPStatusCode: http.StatusNotFound,
		GRPCStatusCode: codes.NotFound,
		PublicMsg:      "Partner profile with the given ID doesn't exist.",
	}
	ErrInvalidOperationType = &AppError{
		Err:            errors.New("invalid operation type"),
		Code:           "INVALID_OPERATION_TYPE",
		HTTPStatusCode: http.StatusBadRequest,
		GRPCStatusCode: codes.InvalidArgument,
		PublicMsg:      "Invalid operation type.",
	}
)
