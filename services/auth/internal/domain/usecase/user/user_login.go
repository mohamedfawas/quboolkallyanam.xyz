package user

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	constants "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	appErrors "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/errors"
	authevents "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/events/auth"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/security/hash"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/auth/internal/domain/entity"
	"gorm.io/gorm"
)

func (u *userUseCase) Login(ctx context.Context, email, password string) (*entity.TokenPair, error) {
	user, err := u.userRepository.GetUser(ctx, "email", email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErrors.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if user == nil {
		return nil, appErrors.ErrUserNotFound
	}

	if !user.IsActive {
		return nil, appErrors.ErrAccountDisabled
	}

	if user.IsBlocked {
		return nil, appErrors.ErrAccountBlocked
	}

	if !hash.VerifyPassword(user.PasswordHash, password) {
		return nil, appErrors.ErrInvalidCredentials
	}

	role := constants.RoleUser
	if user.IsPremium() {
		role = constants.RolePremiumUser
	}

	accessToken, err := u.jwtManager.GenerateAccessToken(user.ID.String(), role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := u.jwtManager.GenerateRefreshToken(user.ID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	refreshTokenKey := fmt.Sprintf("%s%s", constants.RedisPrefixRefreshToken, user.ID.String())
	err = u.tokenRepository.StoreRefreshToken(
		ctx,
		refreshTokenKey,
		refreshToken,
		time.Duration(u.config.Auth.JWT.RefreshTokenDays)*24*time.Hour,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to store user refresh token: %w", err)
	}

	userLoginEvent := authevents.UserLoginSuccessEvent{
		UserID: user.ID,
		Email:  user.Email,
		Phone:  user.Phone,
	}

	if err := u.eventPublisher.PublishUserLoginSuccess(ctx, userLoginEvent); err != nil {
		// No need to fail the login process if the event publishing fails
		log.Printf("failed to publish user login success event: %v", err)
	}

	tokenPair := &entity.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(u.config.Auth.JWT.AccessTokenMinutes) * 60,
	}

	log.Printf("The user is logged in successfully, user id: %s", user.ID.String())

	return tokenPair, nil
}
