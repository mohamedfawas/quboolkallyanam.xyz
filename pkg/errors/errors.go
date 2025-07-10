package errors

import "errors"

// Common errors
var (
	ErrInvalidInput        = errors.New("invalid input parameters")
	ErrInternalServerError = errors.New("internal server error")
	ErrUnauthorized        = errors.New("unauthorized access")
	ErrForbidden           = errors.New("forbidden access")
)

// User authentication errors
var (
	ErrInvalidEmail                = errors.New("invalid email format")
	ErrInvalidPhoneNumber          = errors.New("invalid phone number format")
	ErrInvalidPassword             = errors.New("invalid password format")
	ErrUserNotFound                = errors.New("user not found")
	ErrInvalidCredentials          = errors.New("invalid email or password")
	ErrAccountNotVerified          = errors.New("account is not verified")
	ErrAccountDisabled             = errors.New("account is disabled")
	ErrEmailAlreadyExists          = errors.New("email already registered")
	ErrPhoneAlreadyExists          = errors.New("phone number already registered")
	ErrRegistrationFailed          = errors.New("registration failed")
	ErrVerificationFailed          = errors.New("verification failed")
	ErrHashGenerationFailed        = errors.New("failed to generate hash")
	ErrHashComparisonFailed        = errors.New("failed to compare hash")
	ErrPendingRegistrationNotFound = errors.New("pending registration not found")
)

// Admin authentication errors
var (
	ErrAdminNotFound        = errors.New("admin not found")
	ErrAdminAccountDisabled = errors.New("admin account is disabled")
	ErrInvalidAdminInput    = errors.New("invalid admin input parameters")
)

// Token errors
var (
	ErrInvalidToken        = errors.New("invalid or expired token")
	ErrInvalidRefreshToken = errors.New("invalid or expired refresh token")
	ErrExpiredToken        = errors.New("token has expired")
	ErrTokenNotActive      = errors.New("token not active yet")
)

// OTP errors
var (
	ErrInvalidOTP          = errors.New("invalid or expired OTP")
	ErrOTPGenerationFailed = errors.New("failed to generate OTP")
	ErrOTPNotFound         = errors.New("OTP not found")
)

// SMTP errors
var (
	ErrInvalidConfig = errors.New("invalid SMTP configuration")
)
