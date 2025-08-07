package user

import (
	"context"
	"fmt"
	"time"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apperrors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/auth/internal/domain/entity"
)

func (u *userUseCase) RefreshToken(ctx context.Context, refreshToken string) (*entity.TokenPair, error) {
	claims, err := u.jwtManager.VerifyToken(refreshToken)
	if err != nil {
		return nil, err
	}

	userID := claims.UserID
	refreshTokenKey := fmt.Sprintf("%s%s", constants.RedisPrefixRefreshToken, userID)
	valid, err := u.tokenRepository.IsValidRefreshToken(ctx, refreshTokenKey)
	if err != nil {
		return nil, fmt.Errorf("failed to validate refresh token: %w", err)
	}

	if !valid {
		return nil, apperrors.ErrInvalidToken
	}

	if err := u.tokenRepository.DeleteRefreshToken(ctx, userID); err != nil {
		return nil, fmt.Errorf("failed to delete refresh token: %w", err)
	}

	user, err := u.userRepository.GetUser(ctx, "id", userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if user == nil {
		return nil, apperrors.ErrUserNotFound
	}

	if user.IsBlocked {
		return nil, apperrors.ErrUserBlocked
	}

	role := constants.RoleUser
	if user.IsPremium() {
		role = constants.RolePremiumUser
	}

	newAccessToken, err := u.jwtManager.GenerateAccessToken(userID, role)
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
