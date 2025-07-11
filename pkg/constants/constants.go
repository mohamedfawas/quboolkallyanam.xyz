package constants

const (
	// Commmon
	ID = "id"

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

	// messaging
	UserID           = "user_id"
	EventType        = "event_type"
	EventUserDeleted = "user_deleted"
	Timestamp        = "timestamp"
)
