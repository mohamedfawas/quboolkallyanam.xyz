package constants

const (
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
)
