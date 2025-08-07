package dto

type UserRegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
}

type UserRegisterResponse struct {
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type UserVerificationRequest struct {
	Email string `json:"email" binding:"required,email"`
	OTP   string `json:"otp" binding:"required"`
}

type UserVerificationResponse struct {
	Success bool `json:"success"`
}

type UserLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type UserLoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

type UserLogoutResponse struct {
	Success bool `json:"success"`
}

type AdminLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type AdminLoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

type AdminLogoutResponse struct {
	Success bool `json:"success"`
}

type UserDeleteRequest struct {
	Password string `json:"password" binding:"required,min=8"`
}

type UserDeleteResponse struct {
	Success bool `json:"success"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

type BlockUserRequest struct {
	Field string `json:"field" binding:"required,oneof=email phone id"`
	Value string `json:"value" binding:"required"`
}

type BlockUserResponse struct {
	Success bool `json:"success"`
}

type GetUsersRequest struct {
	Page  int32 `json:"page" binding:"required,min=1"`
	Limit int32 `json:"limit" binding:"required,min=1,max=50"`
}

type GetUserResponse struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	Phone         string `json:"phone"`
	EmailVerified bool   `json:"email_verified"`
	PremiumUntil  *string `json:"premium_until,omitempty"`
	LastLoginAt   *string `json:"last_login_at,omitempty"`
	IsBlocked     bool   `json:"is_blocked"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

type GetUsersResponse struct {
	Users []GetUserResponse `json:"users"`
}


type GetUserByFieldRequest struct {
	Field string `json:"field" binding:"required,oneof=email phone id"`
	Value string `json:"value" binding:"required"`
}

type GetUserByFieldResponse struct {
	User GetUserResponse `json:"user"`
}