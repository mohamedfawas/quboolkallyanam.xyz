package constants

const (
	// Commmon
	ID             = "id"
	EnvProduction  = "production"
	EnvDevelopment = "development"
	EnvStaging     = "staging"

	// Auth
	ContextKeyUserID    = "user_id"
	ContextKeyRole      = "role"
	HeaderAuthorization = "Authorization"
	BearerTokenPrefix   = "Bearer "
	DefaultCostBcrypt   = 12

	// Roles
	RoleUser        = "user"
	RoleAdmin       = "admin"
	RolePremiumUser = "premium_user"

	// Redis key prefixes
	RedisPrefixRefreshToken = "refresh_token:"
	RedisPrefixBlacklist    = "blacklist:"
	RedisPrefixOTP          = "otp:"

	// Token
	BlacklistedToken = "blacklisted"

	// messaging
	UserID           = "user_id"
	EventType        = "event_type"
	EventUserDeleted = "user_deleted"
	Timestamp        = "timestamp"

	// payment
	PaymentMethodRazorpay  = "razorpay"
	PaymentCurrencyINR     = "INR"
	PaymentStatusPending   = "pending"
	PaymentStatusCompleted = "completed"
	PaymentStatusFailed    = "failed"
)
