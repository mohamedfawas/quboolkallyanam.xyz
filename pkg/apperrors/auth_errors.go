package apperrors

import (
	"errors"
	"net/http"

	"google.golang.org/grpc/codes"
)

// common auth errors
var (
	ErrInvalidPassword = &AppError{
		Err:            errors.New("invalid password"),
		Code:           "INVALID_PASSWORD",
		HTTPStatusCode: http.StatusUnauthorized,
		GRPCStatusCode: codes.InvalidArgument,
		PublicMsg:      "Invalid password. Please try again."}
)

// admin errors
var (
	ErrAdminNotFound = &AppError{
		Err:            errors.New("admin not found"),
		Code:           "ADMIN_NOT_FOUND",
		HTTPStatusCode: http.StatusNotFound,
		GRPCStatusCode: codes.NotFound,
		PublicMsg:      "Admin not found. Please try to create your admin account first."}
	ErrAdminAccountDisabled = &AppError{
		Err:            errors.New("admin account disabled"),
		Code:           "ADMIN_ACCOUNT_DISABLED",
		HTTPStatusCode: http.StatusForbidden,
		GRPCStatusCode: codes.PermissionDenied,
		PublicMsg:      "Admin account is disabled. Please contact support."}
	ErrAdminInvalidCredentials = &AppError{
		Err:            errors.New("invalid admin credentials"),
		Code:           "INVALID_ADMIN_CREDENTIALS",
		HTTPStatusCode: http.StatusUnauthorized,
		GRPCStatusCode: codes.InvalidArgument,
		PublicMsg:      "Invalid admin credentials. Please try again."}
)
