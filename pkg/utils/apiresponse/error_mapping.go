package apiresponse

import (
	"net/http"

	appErrors "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ErrorCodeMap maps application errors to standard error codes
var ErrorCodeMap = map[error]string{
	// Auth errors
	appErrors.ErrInvalidCredentials:          "INVALID_CREDENTIALS",
	appErrors.ErrUserNotFound:                "USER_NOT_FOUND",
	appErrors.ErrEmailAlreadyExists:          "EMAIL_ALREADY_EXISTS",
	appErrors.ErrPhoneAlreadyExists:          "PHONE_ALREADY_EXISTS",
	appErrors.ErrAccountNotVerified:          "ACCOUNT_NOT_VERIFIED",
	appErrors.ErrAccountDisabled:             "ACCOUNT_DISABLED",
	appErrors.ErrAccountBlocked:              "ACCOUNT_BLOCKED",
	appErrors.ErrInvalidToken:                "INVALID_TOKEN",
	appErrors.ErrExpiredToken:                "EXPIRED_TOKEN",
	appErrors.ErrInvalidRefreshToken:         "INVALID_REFRESH_TOKEN",
	appErrors.ErrPendingRegistrationNotFound: "PENDING_REGISTRATION_NOT_FOUND",

	// Validation errors
	appErrors.ErrInvalidEmail:       "INVALID_EMAIL",
	appErrors.ErrInvalidPhoneNumber: "INVALID_PHONE",
	appErrors.ErrInvalidPassword:    "INVALID_PASSWORD",
	appErrors.ErrInvalidInput:       "INVALID_INPUT",

	// OTP errors
	appErrors.ErrInvalidOTP:          "INVALID_OTP",
	appErrors.ErrOTPNotFound:         "OTP_NOT_FOUND",
	appErrors.ErrOTPGenerationFailed: "OTP_GENERATION_FAILED",

	// Admin errors
	appErrors.ErrAdminNotFound:        "ADMIN_NOT_FOUND",
	appErrors.ErrAdminAccountDisabled: "ADMIN_ACCOUNT_DISABLED",

	// System errors
	appErrors.ErrInternalServerError: "INTERNAL_SERVER_ERROR",
	appErrors.ErrUnauthorized:        "UNAUTHORIZED",
	appErrors.ErrForbidden:           "FORBIDDEN",
}

// UserFriendlyMessages provides user-friendly error messages
var UserFriendlyMessages = map[error]string{
	appErrors.ErrInvalidCredentials:          "Invalid email or password. Please try again.",
	appErrors.ErrUserNotFound:                "User account not found.",
	appErrors.ErrEmailAlreadyExists:          "An account with this email already exists.",
	appErrors.ErrPhoneAlreadyExists:          "An account with this phone number already exists.",
	appErrors.ErrAccountNotVerified:          "Please verify your account before proceeding.",
	appErrors.ErrAccountDisabled:             "Your account has been disabled. Please contact support.",
	appErrors.ErrAccountBlocked:              "Your account has been blocked. Please contact support.",
	appErrors.ErrInvalidToken:                "Your session has expired. Please log in again.",
	appErrors.ErrExpiredToken:                "Your session has expired. Please log in again.",
	appErrors.ErrInvalidRefreshToken:         "Invalid session. Please log in again.",
	appErrors.ErrPendingRegistrationNotFound: "Registration request not found. Please register again.",
	appErrors.ErrInvalidEmail:                "Please enter a valid email address.",
	appErrors.ErrInvalidPhoneNumber:          "Please enter a valid phone number.",
	appErrors.ErrInvalidPassword:             "Password must be at least 8 characters long with numbers and special characters.",
	appErrors.ErrInvalidOTP:                  "Invalid or expired OTP. Please try again.",
	appErrors.ErrOTPNotFound:                 "OTP not found. Please request a new one.",
	appErrors.ErrAdminNotFound:               "Admin account not found.",
	appErrors.ErrAdminAccountDisabled:        "Admin account is disabled.",
	appErrors.ErrInternalServerError:         "Something went wrong. Please try again later.",
	appErrors.ErrUnauthorized:                "You are not authorized to perform this action.",
	appErrors.ErrForbidden:                   "Access denied.",
}

// mapErrorToResponse maps errors to HTTP status codes and error info
func mapErrorToResponse(err error) (int, *ErrorInfo) {
	// Handle gRPC status errors
	if st, ok := status.FromError(err); ok {
		return mapGRPCStatusToResponse(st)
	}

	// Handle application errors
	if code, exists := ErrorCodeMap[err]; exists {
		message := UserFriendlyMessages[err]
		if message == "" {
			message = err.Error()
		}

		return getHTTPStatusForAppError(err), &ErrorInfo{
			Code:    code,
			Message: message,
		}
	}

	// Handle validation errors (from gin binding)
	if isValidationError(err) {
		return http.StatusBadRequest, &ErrorInfo{
			Code:    "VALIDATION_ERROR",
			Message: "Please check your input and try again.",
			Details: parseValidationError(err),
		}
	}

	// Default to internal server error
	return http.StatusInternalServerError, &ErrorInfo{
		Code:    "INTERNAL_SERVER_ERROR",
		Message: "Something went wrong. Please try again later.",
	}
}

func mapGRPCStatusToResponse(st *status.Status) (int, *ErrorInfo) {
	var httpCode int
	var errorCode string

	switch st.Code() {
	case codes.InvalidArgument:
		httpCode = http.StatusBadRequest
		errorCode = "INVALID_INPUT"
	case codes.Unauthenticated:
		httpCode = http.StatusUnauthorized
		errorCode = "UNAUTHORIZED"
	case codes.PermissionDenied:
		httpCode = http.StatusForbidden
		errorCode = "FORBIDDEN"
	case codes.NotFound:
		httpCode = http.StatusNotFound
		errorCode = "NOT_FOUND"
	case codes.AlreadyExists:
		httpCode = http.StatusConflict
		errorCode = "ALREADY_EXISTS"
	case codes.Unavailable:
		httpCode = http.StatusServiceUnavailable
		errorCode = "SERVICE_UNAVAILABLE"
	default:
		httpCode = http.StatusInternalServerError
		errorCode = "INTERNAL_SERVER_ERROR"
	}

	return httpCode, &ErrorInfo{
		Code:    errorCode,
		Message: st.Message(),
	}
}

func getHTTPStatusForAppError(err error) int {
	switch err {
	case appErrors.ErrInvalidCredentials, appErrors.ErrInvalidToken,
		appErrors.ErrExpiredToken, appErrors.ErrInvalidRefreshToken:
		return http.StatusUnauthorized
	case appErrors.ErrUserNotFound, appErrors.ErrAdminNotFound,
		appErrors.ErrPendingRegistrationNotFound, appErrors.ErrOTPNotFound:
		return http.StatusNotFound
	case appErrors.ErrEmailAlreadyExists, appErrors.ErrPhoneAlreadyExists:
		return http.StatusConflict
	case appErrors.ErrInvalidEmail, appErrors.ErrInvalidPhoneNumber,
		appErrors.ErrInvalidPassword, appErrors.ErrInvalidInput, appErrors.ErrInvalidOTP:
		return http.StatusBadRequest
	case appErrors.ErrAccountDisabled, appErrors.ErrAccountBlocked,
		appErrors.ErrAdminAccountDisabled:
		return http.StatusForbidden
	case appErrors.ErrUnauthorized:
		return http.StatusUnauthorized
	case appErrors.ErrForbidden:
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}

func isValidationError(err error) bool {
	// Add logic to detect validation errors from gin binding
	return false // Implement based on your validation library
}

func parseValidationError(err error) map[string]string {
	// Parse validation errors and return field-specific messages
	return map[string]string{
		"general": err.Error(),
	}
}
