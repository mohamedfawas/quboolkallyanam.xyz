package v1

import (
	authpbv1 "github.com/mohamedfawas/quboolkallyanam.xyz/api/proto/auth/v1"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
	"google.golang.org/protobuf/types/known/emptypb"
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
		Success: resp.Success,
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

func MapUserLogoutResponse(resp *emptypb.Empty) *dto.UserLogoutResponse {
	return &dto.UserLogoutResponse{
		Success: true,
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

func MapAdminLogoutResponse(resp *emptypb.Empty) *dto.AdminLogoutResponse {
	return &dto.AdminLogoutResponse{
		Success: true,
	}
}

////////////////////////////// User Delete //////////////////////////////

func MapUserDeleteRequest(req dto.UserDeleteRequest) *authpbv1.UserDeleteRequest {
	return &authpbv1.UserDeleteRequest{
		Password: req.Password,
	}
}

func MapUserDeleteResponse(resp *emptypb.Empty) *dto.UserDeleteResponse {
	return &dto.UserDeleteResponse{
		Success: true,
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
