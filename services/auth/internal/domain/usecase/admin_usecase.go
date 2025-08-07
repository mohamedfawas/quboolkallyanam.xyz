package usecase

import (
	"context"

	"github.com/mohamedfawas/quboolkallyanam.xyz/services/auth/internal/domain/entity"
)

type AdminUsecase interface {
	InitializeDefaultAdmin(ctx context.Context, defaultEmail, defaultPassword string) error
	AdminLogin(ctx context.Context, email, password string) (*entity.TokenPair, error)
	AdminLogout(ctx context.Context, accessToken string) error
	BlockUser(ctx context.Context, field string, value string) error
	GetUsers(ctx context.Context, page, limit int) ([]*entity.GetUserResponse, error)
	GetUserByField(ctx context.Context, field string, value string) (*entity.GetUserResponse, error)
}
