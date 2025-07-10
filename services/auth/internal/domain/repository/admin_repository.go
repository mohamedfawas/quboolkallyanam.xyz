package repository

import (
	"context"

	"github.com/mohamedfawas/quboolkallyanam.xyz/services/auth/internal/domain/entity"
)

type AdminRepository interface {
	GetAdminByEmail(ctx context.Context, email string) (*entity.Admin, error)
	CreateAdmin(ctx context.Context, admin *entity.Admin) error
	UpdateAdmin(ctx context.Context, admin *entity.Admin) error
	CheckAdminExists(ctx context.Context, email string) (bool, error)
}
