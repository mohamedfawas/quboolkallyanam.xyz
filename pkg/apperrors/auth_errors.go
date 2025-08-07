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
	ErrInvalidEmail = &AppError{
		Err:            errors.New("invalid email"),
		Code:           "INVALID_EMAIL",
		HTTPStatusCode: http.StatusUnauthorized,
		GRPCStatusCode: codes.InvalidArgument,
		PublicMsg:      "Invalid email. Please try again."}
	ErrInvalidPhoneNumber = &AppError{
		Err:            errors.New("invalid phone number"),
		Code:           "INVALID_PHONE_NUMBER",
		HTTPStatusCode: http.StatusUnauthorized,
		GRPCStatusCode: codes.InvalidArgument,
		PublicMsg:      "Invalid phone number. Please try again."}
	ErrOTPNotFound = &AppError{
		Err:            errors.New("otp not found"),
		Code:           "OTP_NOT_FOUND",
		HTTPStatusCode: http.StatusNotFound,
		GRPCStatusCode: codes.NotFound,
		PublicMsg:      "OTP not found. Please request a new OTP."}
	ErrInvalidOTP = &AppError{
		Err:            errors.New("invalid otp"),
		Code:           "INVALID_OTP",
		HTTPStatusCode: http.StatusUnauthorized,
		GRPCStatusCode: codes.InvalidArgument,
		PublicMsg:      "Invalid OTP. Please request a new OTP."}
	ErrInvalidCredentials = &AppError{
		Err:            errors.New("invalid credentials"),
		Code:           "INVALID_CREDENTIALS",
		HTTPStatusCode: http.StatusUnauthorized,
		GRPCStatusCode: codes.InvalidArgument,
		PublicMsg:      "Invalid credentials. Please try again."}
)

// pending registration errors
var (
	ErrPendingRegistrationNotFound = &AppError{
		Err:            errors.New("pending registration not found"),
		Code:           "PENDING_REGISTRATION_NOT_FOUND",
		HTTPStatusCode: http.StatusNotFound,
		GRPCStatusCode: codes.NotFound,
		PublicMsg:      "Pending registration not found. Please try again."}
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

// user errors
var (
	ErrUserAlreadyBlocked = &AppError{
		Err:            errors.New("user already blocked"),
		Code:           "USER_ALREADY_BLOCKED",
		HTTPStatusCode: http.StatusBadRequest,
		GRPCStatusCode: codes.InvalidArgument,
		PublicMsg:      "User already blocked."}
	ErrEmailAlreadyExists = &AppError{
		Err:            errors.New("email already exists"),
		Code:           "EMAIL_ALREADY_EXISTS",
		HTTPStatusCode: http.StatusBadRequest,
		GRPCStatusCode: codes.AlreadyExists,
		PublicMsg:      "Email already exists. Please try with a different email."}
	ErrPhoneAlreadyExists = &AppError{
		Err:            errors.New("phone already exists"),
		Code:           "PHONE_ALREADY_EXISTS",
		HTTPStatusCode: http.StatusBadRequest,
		GRPCStatusCode: codes.AlreadyExists,
		PublicMsg:      "Phone already exists. Please try with a different phone number."}
	ErrUserBlocked = &AppError{
		Err:            errors.New("user blocked"),
		Code:           "USER_BLOCKED",
		HTTPStatusCode: http.StatusForbidden,
		GRPCStatusCode: codes.PermissionDenied,
		PublicMsg:      "User blocked. Please contact support."}
)

// Token errors
var (
	ErrInvalidToken = &AppError{
		Err:            errors.New("invalid token"),
		Code:           "INVALID_TOKEN",
		HTTPStatusCode: http.StatusUnauthorized,
		GRPCStatusCode: codes.InvalidArgument,
		PublicMsg:      "Invalid token. Please try again."}
	ErrExpiredToken = &AppError{
		Err:            errors.New("expired token"),
		Code:           "EXPIRED_TOKEN",
		HTTPStatusCode: http.StatusUnauthorized,
		GRPCStatusCode: codes.InvalidArgument,
		PublicMsg:      "Expired token. Please try again."}
	ErrTokenNotActive = &AppError{
		Err:            errors.New("token not active"),
		Code:           "TOKEN_NOT_ACTIVE",
		HTTPStatusCode: http.StatusUnauthorized,
		GRPCStatusCode: codes.InvalidArgument,
		PublicMsg:      "Token not active. Please try again."}
	ErrUnauthorized = &AppError{
		Err:            errors.New("unauthorized"),
		Code:           "UNAUTHORIZED",
		HTTPStatusCode: http.StatusUnauthorized,
		GRPCStatusCode: codes.Unauthenticated,
		PublicMsg:      "Unauthorized. Please try again."}
	ErrForbidden = &AppError{
		Err:            errors.New("forbidden"),
		Code:           "FORBIDDEN",
		HTTPStatusCode: http.StatusForbidden,
		GRPCStatusCode: codes.PermissionDenied,
		PublicMsg:      "Forbidden. Please try again."}
	ErrAccessTokenNotFound = &AppError{
		Err:            errors.New("access token not found"),
		Code:           "ACCESS_TOKEN_NOT_FOUND",
		HTTPStatusCode: http.StatusUnauthorized,
		GRPCStatusCode: codes.Unauthenticated,
		PublicMsg:      "Access token not found. Please try again by giving a valid access token."}
	ErrRefreshTokenNotFound = &AppError{
		Err:            errors.New("refresh token not found"),
		Code:           "REFRESH_TOKEN_NOT_FOUND",
		HTTPStatusCode: http.StatusUnauthorized,
		GRPCStatusCode: codes.Unauthenticated,
		PublicMsg:      "Refresh token not found. Please try again by giving a valid refresh token."}
)
