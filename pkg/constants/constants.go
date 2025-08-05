package constants

const (
	// Commmon
	ID              = "id"
	EnvProduction   = "production"
	EnvDevelopment  = "development"
	EnvStaging      = "staging"
	ENVIRONMENT     = "ENVIRONMENT"
	Port            = "port"
	Service         = "service"
	ServiceGateway  = "gateway"
	ServiceAuth     = "auth"
	ServiceUser     = "user"
	ServiceAdmin    = "admin"
	ServicePayment  = "payment"
	ServiceChat     = "chat"
	ServiceNotification = "notification"
	Env             = "env"
	HeaderRequestID = "X-Request-ID"
	Unknown         = "unknown"
	HTTPStatusCode  = "http_status_code"
	UserFriendlyMessage     = "user_message"

	InteralServerErrorMessage = "Something went wrong. Please try again later."

	// Operation Types
	CreateOperationType = "create"
	UpdateOperationType = "update"

	// HTTP Handler Versions
	HTTPHandlerVersionV1 = "/api/v1"

	// Auth
	ContextKeyUserID      = "user_id"
	ContextKeyRole        = "role"
	ContextKeyRequestID   = "request_id"
	ContextKeyAccessToken = "access_token"
	HeaderAuthorization   = "Authorization"
	HeaderRefreshToken    = "Refresh-Token"
	BearerTokenPrefix     = "Bearer "
	DefaultCostBcrypt     = 12

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
	PaymentStatusExpired   = "expired"

	// subscription
	SubscriptionStatusActive    = "active"
	SubscriptionStatusCancelled = "cancelled"

	// MongoDB Collections
	MongoDBCollectionConversations = "conversations"
	MongoDBCollectionMessages      = "messages"

	// Conversation display limits
	DefaultConversationDisplayLimit = 10
	MaxConversationDisplayLimit     = 50

	DefaultPaginationLimit = 10
	MaxPaginationLimit = 50

	// Image file
	ImageFileMaxSize = 5 * 1024 * 1024

	// gcs storage
	ProfilePhotoStorageDirectory = "profile-photos"
	AdditionalPhotoStorageDirectory = "additional-photos"
)
