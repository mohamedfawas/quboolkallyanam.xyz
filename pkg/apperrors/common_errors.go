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
	ErrInvalidField = &AppError{
		Err:            errors.New("invalid field"),
		Code:           "INVALID_FIELD",
		HTTPStatusCode: http.StatusBadRequest,
		GRPCStatusCode: codes.InvalidArgument,
		PublicMsg:      "The field you requested to update is non existent.",
	}
	ErrInvalidPaginationPage = &AppError{
		Err:            errors.New("invalid pagination page"),
		Code:           "INVALID_PAGINATION_PAGE",
		HTTPStatusCode: http.StatusBadRequest,
		GRPCStatusCode: codes.InvalidArgument,
		PublicMsg:      "Invalid pagination page.",
	}
	ErrInvalidPaginationLimit = &AppError{		
		Err:            errors.New("invalid pagination limit"),
		Code:           "INVALID_PAGINATION_LIMIT",
		HTTPStatusCode: http.StatusBadRequest,
		GRPCStatusCode: codes.InvalidArgument,
		PublicMsg:      "Invalid pagination limit.",
	}
	ErrMissingRequiredFields = &AppError{
		Err:            errors.New("missing required fields"),
		Code:           "MISSING_REQUIRED_FIELDS",
		HTTPStatusCode: http.StatusBadRequest,
		GRPCStatusCode: codes.InvalidArgument,
		PublicMsg:      "Missing required fields.",
	}
	ErrInvalidInput = &AppError{
		Err:            errors.New("invalid input"),
		Code:           "INVALID_INPUT",
		HTTPStatusCode: http.StatusBadRequest,
		GRPCStatusCode: codes.InvalidArgument,
		PublicMsg:      "Invalid input.",
	}
)
