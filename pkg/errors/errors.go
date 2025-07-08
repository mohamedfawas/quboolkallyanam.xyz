package errors

import "errors"

// Internal server error
var (
	ErrInternalServerError = errors.New("internal server error")
)

// User authentication errors
var (
	ErrInvalidEmail       = errors.New("invalid email format")
	ErrInvalidPhoneNumber = errors.New("invalid phone number format")
	ErrInvalidPassword    = errors.New("invalid password format")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrAccountNotVerified = errors.New("account is not verified")
	ErrAccountDisabled    = errors.New("account is disabled")
	ErrEmailAlreadyExists = errors.New("email already registered")
	ErrPhoneAlreadyExists = errors.New("phone number already registered")
	ErrRegistrationFailed = errors.New("registration failed")
	ErrVerificationFailed = errors.New("verification failed")
)

// Token errors
var (
	ErrInvalidToken        = errors.New("invalid or expired token")
	ErrInvalidRefreshToken = errors.New("invalid or expired refresh token")
)

// OTP errors
var (
	ErrInvalidOTP          = errors.New("invalid or expired OTP")
	ErrOTPGenerationFailed = errors.New("failed to generate OTP")
)
