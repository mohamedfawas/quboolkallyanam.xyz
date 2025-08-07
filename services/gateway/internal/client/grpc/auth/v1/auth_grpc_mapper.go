package v1

import (
	authpbv1 "github.com/mohamedfawas/quboolkallyanam.xyz/api/proto/auth/v1"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
)

////////////////////////////// User Registration //////////////////////////////

func MapUserRegisterRequest(req dto.UserRegisterRequest) *authpbv1.UserRegisterRequest {
	return &authpbv1.UserRegisterRequest{
		Email:    req.Email,
		Phone:    req.Phone,
		Password: req.Password,
	}
}

func MapUserRegisterResponse(resp *authpbv1.UserRegisterResponse) *dto.UserRegisterResponse {
	return &dto.UserRegisterResponse{
		Email: resp.Email,
		Phone: resp.Phone,
	}
}

////////////////////////////// User Verification //////////////////////////////

func MapUserVerificationRequest(req dto.UserVerificationRequest) *authpbv1.UserVerificationRequest {
	return &authpbv1.UserVerificationRequest{
		Email: req.Email,
		Otp:   req.OTP,
	}
}

func MapUserVerificationResponse(resp *authpbv1.UserVerificationResponse) *dto.UserVerificationResponse {
	return &dto.UserVerificationResponse{
		Success: resp.Success.GetValue(),
	}
}

////////////////////////////// User Login //////////////////////////////

func MapUserLoginRequest(req dto.UserLoginRequest) *authpbv1.UserLoginRequest {
	return &authpbv1.UserLoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}
}

func MapUserLoginResponse(resp *authpbv1.UserLoginResponse) *dto.UserLoginResponse {
	return &dto.UserLoginResponse{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
		ExpiresIn:    resp.ExpiresIn,
	}
}

////////////////////////////// User Logout //////////////////////////////

func MapUserLogoutRequest(accessToken string) *authpbv1.UserLogoutRequest {
	return &authpbv1.UserLogoutRequest{
		AccessToken: accessToken,
	}
}

func MapUserLogoutResponse(resp *authpbv1.UserLogoutResponse) *dto.UserLogoutResponse {
	return &dto.UserLogoutResponse{
		Success: resp.Success.GetValue(),
	}
}

////////////////////////////// Admin Login //////////////////////////////

func MapAdminLoginRequest(req dto.AdminLoginRequest) *authpbv1.AdminLoginRequest {
	return &authpbv1.AdminLoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}
}

func MapAdminLoginResponse(resp *authpbv1.AdminLoginResponse) *dto.AdminLoginResponse {
	return &dto.AdminLoginResponse{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
		ExpiresIn:    resp.ExpiresIn,
	}
}

////////////////////////////// Admin Logout //////////////////////////////

func MapAdminLogoutRequest(accessToken string) *authpbv1.AdminLogoutRequest {
	return &authpbv1.AdminLogoutRequest{
		AccessToken: accessToken,
	}
}

func MapAdminLogoutResponse(resp *authpbv1.AdminLogoutResponse) *dto.AdminLogoutResponse {
	return &dto.AdminLogoutResponse{
		Success: resp.Success.GetValue(),
	}
}

////////////////////////////// User Delete //////////////////////////////

func MapUserDeleteRequest(req dto.UserDeleteRequest) *authpbv1.UserDeleteRequest {
	return &authpbv1.UserDeleteRequest{
		Password: req.Password,
	}
}

func MapUserDeleteResponse(resp *authpbv1.UserDeleteResponse) *dto.UserDeleteResponse {
	return &dto.UserDeleteResponse{
		Success: resp.Success.GetValue(),
	}
}

// //////////////////////////// Refresh Token //////////////////////////////
func MapRefreshTokenRequest(req dto.RefreshTokenRequest) *authpbv1.RefreshTokenRequest {
	return &authpbv1.RefreshTokenRequest{
		RefreshToken: req.RefreshToken,
	}
}

func MapRefreshTokenResponse(resp *authpbv1.RefreshTokenResponse) *dto.RefreshTokenResponse {
	return &dto.RefreshTokenResponse{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
		ExpiresIn:    resp.ExpiresIn,
	}
}

/////////////////////////// Block user //////////////////////////////
func MapBlockUserRequest(req dto.BlockUserRequest) *authpbv1.BlockUserRequest {
	return &authpbv1.BlockUserRequest{
		Field: req.Field,
		Value: req.Value,
	}
}

func MapBlockUserResponse(resp *authpbv1.BlockUserResponse) *dto.BlockUserResponse {
	return &dto.BlockUserResponse{
		Success: resp.Success.GetValue(),
	}
}

/////////////////////////// Get Users //////////////////////////////
func MapGetUsersRequest(req dto.GetUsersRequest) *authpbv1.GetUsersRequest {
	return &authpbv1.GetUsersRequest{
		Page:  req.Page,
		Limit: req.Limit,
	}
}

func MapGetUsersResponse(resp *authpbv1.GetUsersResponse) *dto.GetUsersResponse {
	users := make([]dto.GetUserResponse, len(resp.Users))
	for i, user := range resp.Users {
		users[i] = MapGetUserResponse(user)
	}
	return &dto.GetUsersResponse{
		Users: users,
	}
}

/////////////////////////// Get User By Field //////////////////////////////
func MapGetUserByFieldRequest(req dto.GetUserByFieldRequest) *authpbv1.GetUserByFieldRequest {
	return &authpbv1.GetUserByFieldRequest{
		Field: req.Field,
		Value: req.Value,
	}
}

func MapGetUserByFieldResponse(resp *authpbv1.GetUserByFieldResponse) *dto.GetUserByFieldResponse {
	return &dto.GetUserByFieldResponse{
		User: MapGetUserResponse(resp.User),
	}
}

/////////////////////////// Get User Response Helper //////////////////////////////
func MapGetUserResponse(user *authpbv1.GetUserResponse) dto.GetUserResponse {
	var premiumUntil *string
	if user.PremiumUntil != nil {
		formatted := user.PremiumUntil.AsTime().Format("2006-01-02T15:04:05Z07:00")
		premiumUntil = &formatted
	}

	var lastLoginAt *string
	if user.LastLoginAt != nil {
		formatted := user.LastLoginAt.AsTime().Format("2006-01-02T15:04:05Z07:00")
		lastLoginAt = &formatted
	}

	return dto.GetUserResponse{
		ID:            user.Id,
		Email:         user.Email,
		Phone:         user.Phone,
		EmailVerified: user.EmailVerified,
		PremiumUntil:  premiumUntil,
		LastLoginAt:   lastLoginAt,
		IsBlocked:     user.IsBlocked,
		CreatedAt:     user.CreatedAt.AsTime().Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:     user.UpdatedAt.AsTime().Format("2006-01-02T15:04:05Z07:00"),
	}
}