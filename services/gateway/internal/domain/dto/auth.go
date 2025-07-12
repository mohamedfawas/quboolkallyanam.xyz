package dto

type UserRegisterRequest struct {
	Email    string `json:"email" binding:"required"`
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserRegisterResponse struct {
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type UserVerificationRequest struct {
	Email string `json:"email" binding:"required"`
	OTP   string `json:"otp" binding:"required"`
}

type UserVerificationResponse struct {
	Success string `json:"success"`
}

type UserLoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserLoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

type UserLogoutRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type UserLogoutResponse struct {
	Success bool `json:"success"`
}

type AdminLoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AdminLoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

type AdminLogoutRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type AdminLogoutResponse struct {
	Success bool `json:"success"`
}

type UserDeleteRequest struct {
	Password string `json:"password" binding:"required"`
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
