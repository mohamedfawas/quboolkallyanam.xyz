package usecase

import (
	"context"

	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/config"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
)

type AuthUsecase interface {
	UserRegister(ctx context.Context, req dto.UserRegisterRequest) (*dto.UserRegisterResponse, error)
	UserVerification(ctx context.Context, req dto.UserVerificationRequest, config config.Config) (*dto.UserVerificationResponse, error)
	UserLogin(ctx context.Context, req dto.UserLoginRequest) (*dto.UserLoginResponse, error)
	UserLogout(ctx context.Context, accessToken string) error
	UserDelete(ctx context.Context, req dto.UserDeleteRequest) error
	AdminLogin(ctx context.Context, req dto.AdminLoginRequest) (*dto.AdminLoginResponse, error)
	AdminLogout(ctx context.Context, accessToken string) error
	RefreshToken(ctx context.Context, req dto.RefreshTokenRequest) (*dto.RefreshTokenResponse, error)
	BlockOrUnblockUser(ctx context.Context, req dto.BlockOrUnblockUserRequest) (*dto.BlockOrUnblockUserResponse, error)
	GetUsers(ctx context.Context, req dto.GetUsersRequest) (*dto.GetUsersResponse, error)
	GetUserByField(ctx context.Context, req dto.GetUserByFieldRequest) (*dto.GetUserByFieldResponse, error)
}
