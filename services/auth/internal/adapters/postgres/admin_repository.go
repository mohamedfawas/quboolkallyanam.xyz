package postgres

import (
	"context"
	"errors"

	"gorm.io/gorm"

	postgres "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/database/postgres"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/auth/internal/domain/entity"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/auth/internal/domain/repository"
)

type adminRepository struct {
	db *postgres.Client
}

func NewAdminRepository(db *postgres.Client) repository.AdminRepository {
	return &adminRepository{db: db}
}

func (r *adminRepository) GetAdminByEmail(
	ctx context.Context, 
	email string) (*entity.Admin, error) {
	
	var admin entity.Admin
	if err := r.db.GormDB.WithContext(ctx).Where("email = ?", email).First(&admin).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &admin, nil
}

func (r *adminRepository) CreateAdmin(ctx context.Context, admin *entity.Admin) error {
	return r.db.GormDB.WithContext(ctx).Create(admin).Error
}

func (r *adminRepository) UpdateAdmin(ctx context.Context, admin *entity.Admin) error {
	return r.db.GormDB.WithContext(ctx).Save(admin).Error
}

func (r *adminRepository) CheckAdminExists(ctx context.Context, email string) (bool, error) {
	var count int64
	if err := r.db.GormDB.WithContext(ctx).Model(&entity.Admin{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
