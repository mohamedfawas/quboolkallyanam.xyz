package usecase

import (
	"context"

	"github.com/mohamedfawas/quboolkallyanam.xyz/services/auth/internal/domain/entity"
)

type AdminUsecase interface {
	InitializeDefaultAdmin(ctx context.Context, defaultEmail, defaultPassword string) error
	AdminLogin(ctx context.Context, email, password string) (*entity.TokenPair, error)
	AdminLogout(ctx context.Context, refreshToken string) error
}
