package user

import (
	"context"
	"fmt"
	"time"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apperrors"
	constants "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
)

func (u *userUseCase) Logout(ctx context.Context, accessToken string) error {
	claims, err := u.jwtManager.VerifyToken(accessToken)
	if err != nil {
		return err
	}

	if claims.Role != constants.RoleUser {
		return apperrors.ErrUnauthorized
	}

	tokenID := claims.ID
	blacklistKey := fmt.Sprintf("%s%s", constants.RedisPrefixBlacklist, tokenID)
	isBlacklisted, err := u.tokenRepository.IsTokenBlacklisted(ctx, blacklistKey)
	if err != nil {
		return fmt.Errorf("failed to check if token is blacklisted: %w", err)
	}

	if isBlacklisted {
		return apperrors.ErrInvalidToken
	}

	timeUntilExpiry := time.Until(claims.ExpiresAt.Time)
	if timeUntilExpiry <= 0 {
		return apperrors.ErrExpiredToken
	}

	err = u.tokenRepository.BlacklistToken(ctx, blacklistKey, timeUntilExpiry)
	if err != nil {
		return fmt.Errorf("failed to blacklist token: %w", err)
	}

	userID := claims.UserID
	refreshTokenKey := fmt.Sprintf("%s%s", constants.RedisPrefixRefreshToken, userID)
	err = u.tokenRepository.DeleteRefreshToken(ctx, refreshTokenKey)
	if err != nil {
		return fmt.Errorf("failed to delete refresh token: %w", err)
	}

	return nil
}
