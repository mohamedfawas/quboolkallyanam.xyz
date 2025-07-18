package user

import (
	"context"
	"fmt"
	"time"

	constants "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	appErrors "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/errors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/auth/internal/domain/entity"
)

func (u *userUseCase) RefreshToken(ctx context.Context, refreshToken string) (*entity.TokenPair, error) {
	claims, err := u.jwtManager.VerifyToken(refreshToken)
	if err != nil {
		return nil, appErrors.ErrInvalidToken
	}

	userID := claims.UserID
	refreshTokenKey := fmt.Sprintf("%s%s", constants.RedisPrefixRefreshToken, userID)
	valid, err := u.tokenRepository.IsValidRefreshToken(ctx, refreshTokenKey)
	if err != nil {
		return nil, fmt.Errorf("failed to validate refresh token: %w", err)
	}

	if !valid {
		return nil, appErrors.ErrInvalidToken
	}

	if err := u.tokenRepository.DeleteRefreshToken(ctx, userID); err != nil {
		return nil, fmt.Errorf("failed to delete refresh token: %w", err)
	}

	newAccessToken, err := u.jwtManager.GenerateAccessToken(userID, constants.RoleUser)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	newRefreshToken, err := u.jwtManager.GenerateRefreshToken(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	refreshTokenKey = fmt.Sprintf("%s%s", constants.RedisPrefixRefreshToken, userID)
	err = u.tokenRepository.StoreRefreshToken(
		ctx,
		refreshTokenKey,
		newRefreshToken,
		time.Duration(u.config.Auth.JWT.RefreshTokenDays)*24*time.Hour,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	newTokenPair := &entity.TokenPair{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    int64(u.config.Auth.JWT.AccessTokenMinutes) * 60,
	}

	return newTokenPair, nil
}
