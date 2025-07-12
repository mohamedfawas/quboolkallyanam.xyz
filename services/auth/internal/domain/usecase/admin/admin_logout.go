package admin

import (
	"context"
	"fmt"
	"time"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	appErrors "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/errors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/logger"
)

func (u *adminUsecase) AdminLogout(ctx context.Context, accessToken string) error {
	claims, err := u.jwtManager.VerifyToken(accessToken)
	if err != nil {
		switch err {
		case appErrors.ErrExpiredToken:
			return appErrors.ErrExpiredToken
		case appErrors.ErrTokenNotActive:
			return appErrors.ErrTokenNotActive
		default:
			return appErrors.ErrInvalidToken
		}
	}

	if claims.Role != constants.RoleAdmin {
		return appErrors.ErrUnauthorized
	}

	tokenID := claims.ID
	blacklistKey := fmt.Sprintf("%s%s", constants.RedisPrefixBlacklist, tokenID)
	isBlacklisted, err := u.tokenRepository.IsTokenBlacklisted(ctx, blacklistKey)
	if err != nil {
		return fmt.Errorf("failed to check if token is blacklisted: %w", err)
	}

	if isBlacklisted {
		return appErrors.ErrInvalidToken
	}

	timeUntilExpiry := time.Until(claims.ExpiresAt.Time)
	if timeUntilExpiry <= 0 {
		return appErrors.ErrExpiredToken
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

	logger.Log.Info("Admin logged out successfully : ", "admin_id", userID)
	return nil
}
