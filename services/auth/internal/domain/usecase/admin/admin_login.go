package admin

import (
	"context"
	"fmt"
	"time"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apperrors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/security/hash"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/auth/internal/domain/entity"
)

func (u *adminUsecase) AdminLogin(ctx context.Context, email, password string) (*entity.TokenPair, error) {
	admin, err := u.adminRepository.GetAdminByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to get admin details using email: %w", err)
	}

	if admin == nil {
		return nil, apperrors.ErrAdminNotFound
	}

	if !admin.IsActive {
		return nil, apperrors.ErrAdminAccountDisabled
	}

	if !hash.VerifyPassword(admin.PasswordHash, password) {
		return nil, apperrors.ErrAdminInvalidCredentials
	}

	accessToken, err := u.jwtManager.GenerateAccessToken(admin.ID.String(), constants.RoleAdmin)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := u.jwtManager.GenerateRefreshToken(admin.ID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	refreshTokenKey := fmt.Sprintf("%s%s", constants.RedisPrefixRefreshToken, admin.ID.String())
	err = u.tokenRepository.StoreRefreshToken(
		ctx,
		refreshTokenKey,
		refreshToken,
		time.Duration(u.config.Auth.JWT.RefreshTokenDays)*24*time.Hour,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to store admin refresh token: %w", err)
	}

	tokenPair := &entity.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(u.config.Auth.JWT.AccessTokenMinutes) * 60,
	}

	return tokenPair, nil
}
