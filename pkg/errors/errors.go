package errors

import "errors"

// Common errors
var (
	ErrInvalidInput         = errors.New("invalid input parameters")
	ErrInternalServerError  = errors.New("internal server error")
	ErrUnauthorized         = errors.New("unauthorized access")
	ErrForbidden            = errors.New("forbidden access")
	ErrInvalidOperationType = errors.New("invalid operation type")
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
	ErrAccountBlocked              = errors.New("account is blocked")
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

// Subscription Plan errors
var (
	ErrSubscriptionPlanNotFound  = errors.New("subscription plan not found")
	ErrSubscriptionPlanNotActive = errors.New("subscription plan is not active")
)

// Payment errors
var (
	ErrPaymentNotFound         = errors.New("payment not found")
	ErrPaymentAlreadyCompleted = errors.New("payment already completed")
	ErrPaymentExpired          = errors.New("payment has expired")
	ErrPaymentSignatureInvalid = errors.New("payment signature verification failed")
	ErrRazorpayOrderCreation   = errors.New("failed to create payment order")
	ErrPaymentProcessingFailed = errors.New("payment processing failed")
)

// User errors

var (
	ErrInvalidFullName       = errors.New("invalid full name: must contain only English alphabets and be 2 to 100 characters long")
	ErrInvalidDateOfBirth    = errors.New("invalid date of birth: must be in the format YYYY-MM-DD and be less than 100 years old")
	ErrInvalidHeight         = errors.New("invalid height: must be between 100 and 250 cm")
	ErrInvalidCommunity      = errors.New("invalid community")
	ErrInvalidMaritalStatus  = errors.New("invalid marital status")
	ErrInvalidProfession     = errors.New("invalid profession")
	ErrInvalidProfessionType = errors.New("invalid profession type")
	ErrInvalidEducationLevel = errors.New("invalid education level")
	ErrInvalidHomeDistrict   = errors.New("invalid home district")

	// Partner Preference errors
	ErrInvalidAgeRange    = errors.New("invalid age range: must be between 18 and 100 years")
	ErrInvalidHeightRange = errors.New("invalid height range: must be between 130 and 220 cm")
	ErrPartnerPreferencesAlreadyExists = errors.New("partner preferences already exists")
	ErrPartnerPreferencesNotFound = errors.New("partner preferences not found")

	ErrUserProfileNotFound = errors.New("user profile not found")

	// Match Making
	ErrInvalidMatchAction = errors.New("invalid match action")
	ErrInvalidTargetProfileID = errors.New("invlaid target profile id")
)
