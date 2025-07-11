package usecase

import (
	"context"

	"github.com/mohamedfawas/quboolkallyanam.xyz/services/auth/internal/domain/entity"
)

type UserUsecase interface {
	Login(ctx context.Context, email, password string) (*entity.TokenPair, error)
	Logout(ctx context.Context, accessToken string) error
}
