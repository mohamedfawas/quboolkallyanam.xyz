package admin

import (
	"context"
	"fmt"
	"time"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	errors "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/errors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/logger"
)

func (u *adminUsecase) AdminLogout(ctx context.Context, accessToken string) error {
	claims, err := u.jwtManager.VerifyToken(accessToken)
	if err != nil {
		switch err {
		case errors.ErrExpiredToken:
			return errors.ErrExpiredToken
		case errors.ErrTokenNotActive:
			return errors.ErrTokenNotActive
		default:
			return errors.ErrInvalidToken
		}
	}

	if claims.Role != constants.RoleAdmin {
		return errors.ErrUnauthorized
	}

	tokenID := claims.ID
	isBlacklisted, err := u.tokenRepository.IsTokenBlacklisted(ctx, tokenID)
	if err != nil {
		return fmt.Errorf("failed to check if token is blacklisted: %w", err)
	}

	if isBlacklisted {
		return errors.ErrInvalidToken
	}

	timeUntilExpiry := time.Until(claims.ExpiresAt.Time)
	if timeUntilExpiry <= 0 {
		return errors.ErrExpiredToken
	}
	err = u.tokenRepository.BlacklistToken(ctx, tokenID, timeUntilExpiry)
	if err != nil {
		return fmt.Errorf("failed to blacklist token: %w", err)
	}

	userID := claims.UserID
	err = u.tokenRepository.DeleteRefreshToken(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to delete refresh token: %w", err)
	}

	logger.Log.Info("Admin logged out successfully : ", "admin_id", userID)
	return nil
}
