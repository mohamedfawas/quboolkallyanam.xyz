package usecase

import (
	"context"
	"time"

	"github.com/mohamedfawas/quboolkallyanam.xyz/services/auth/internal/domain/entity"
)

type UserUsecase interface {
	Login(ctx context.Context, email, password string) (*entity.TokenPair, error)
	Logout(ctx context.Context, accessToken string) error
	RefreshToken(ctx context.Context, refreshToken string) (*entity.TokenPair, error)
	UserAccountDelete(ctx context.Context, userID string, password string) error
	UpdateUserPremium(ctx context.Context, userID string, premiumUntil time.Time) error
}
