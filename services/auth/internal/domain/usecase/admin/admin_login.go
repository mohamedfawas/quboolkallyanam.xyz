package admin

import (
	"context"
	"fmt"
	"time"

	constants "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	errors "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/errors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/security/hash"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/auth/internal/domain/entity"
)

func (u *adminUsecase) AdminLogin(ctx context.Context, email, password string) (*entity.TokenPair, error) {
	admin, err := u.adminRepository.GetAdminByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to get admin details using email: %w", err)
	}

	if admin == nil {
		return nil, errors.ErrAdminNotFound
	}

	if !admin.IsActive {
		return nil, errors.ErrAdminAccountDisabled
	}

	if !hash.ComparePassword(password, admin.PasswordHash) {
		return nil, errors.ErrInvalidPassword
	}

	accessToken, err := u.jwtManager.GenerateAccessToken(admin.ID.String(), constants.RoleAdmin)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := u.jwtManager.GenerateRefreshToken(admin.ID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	err = u.tokenRepository.StoreRefreshToken(
		ctx,
		admin.ID.String(),
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
