package postgres

import (
	"context"
	"log"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/database/postgres"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/auth/internal/domain/entity"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/auth/internal/domain/repository"
)

type pendingRegistrationRepository struct {
	db *postgres.Client
}

func NewPendingRegistrationRepository(db *postgres.Client) repository.PendingRegistrationRepository {
	return &pendingRegistrationRepository{db: db}
}

func (r *pendingRegistrationRepository) CreatePendingRegistration(ctx context.Context, pendingRegistration *entity.PendingRegistration) error {
	return r.db.GormDB.WithContext(ctx).Create(pendingRegistration).Error
}

func (r *pendingRegistrationRepository) GetPendingRegistration(ctx context.Context, field, value string) (*entity.PendingRegistration, error) {
	var pendingRegistration entity.PendingRegistration
	if err := r.db.GormDB.WithContext(ctx).Where(field+" = ?", value).First(&pendingRegistration).Error; err != nil {
		log.Printf("GetPendingRegistration error in pending registration repository: %v", err)
		return nil, err
	}
	return &pendingRegistration, nil
}

func (r *pendingRegistrationRepository) DeletePendingRegistration(ctx context.Context, id int) error {
	return r.db.GormDB.WithContext(ctx).Delete(&entity.PendingRegistration{}, id).Error
}
