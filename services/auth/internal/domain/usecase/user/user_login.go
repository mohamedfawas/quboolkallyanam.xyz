package user

import (
	"context"
	"fmt"
	"time"

	constants "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	errors "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/errors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/security/hash"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/auth/internal/domain/entity"
)

func (u *userUseCase) Login(ctx context.Context, email, password string) (*entity.TokenPair, error) {
	user, err := u.userRepository.GetUser(ctx, "email", email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if user == nil {
		return nil, errors.ErrUserNotFound
	}

	if !user.IsActive {
		return nil, errors.ErrAccountDisabled
	}

	if user.IsBlocked {
		return nil, errors.ErrAccountBlocked
	}

	if hash.ComparePassword(user.PasswordHash, password) {
		return nil, errors.ErrInvalidCredentials
	}

	accessToken, err := u.jwtManager.GenerateAccessToken(user.ID.String(), constants.RoleUser)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := u.jwtManager.GenerateRefreshToken(user.ID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	err = u.tokenRepository.StoreRefreshToken(
		ctx,
		user.ID.String(),
		refreshToken,
		time.Duration(u.config.Auth.JWT.RefreshTokenDays)*24*time.Hour,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to store user refresh token: %w", err)
	}

	tokenPair := &entity.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(u.config.Auth.JWT.AccessTokenMinutes) * 60,
	}

	return tokenPair, nil
}
