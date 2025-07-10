package repository

import (
	"context"

	"github.com/mohamedfawas/quboolkallyanam.xyz/services/auth/internal/domain/entity"
)

type PendingRegistrationRepository interface {
	CreatePendingRegistration(ctx context.Context, registration *entity.PendingRegistration) error
	GetPendingRegistration(ctx context.Context, field, value string) (*entity.PendingRegistration, error)
	DeletePendingRegistration(ctx context.Context, id int) error
}
