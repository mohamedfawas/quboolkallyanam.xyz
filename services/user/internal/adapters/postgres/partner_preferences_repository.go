package postgres


import (
	"context"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/database/postgres"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/user/internal/domain/entity"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/user/internal/domain/repository"
	"gorm.io/gorm"
)

type partnerPreferencesRepository struct {
	db *postgres.Client
}

func NewPartnerPreferencesRepository(db *postgres.Client) repository.PartnerPreferencesRepository {
	return &partnerPreferencesRepository{db: db}
}

func (r *partnerPreferencesRepository) CreatePartnerPreferences(
	ctx context.Context,
	preferences *entity.PartnerPreference) error {
	
	return r.db.GormDB.WithContext(ctx).Create(preferences).Error
}

func (r *partnerPreferencesRepository) GetPartnerPreferencesByUserProfileID(
	ctx context.Context, 
	userProfileID int64) (*entity.PartnerPreference, error) {

	var preferences entity.PartnerPreference
	err := r.db.GormDB.WithContext(ctx).
		Where("user_profile_id = ?", userProfileID).
		First(&preferences).Error
	
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	
	return &preferences, nil
}


func (r *partnerPreferencesRepository) PatchPartnerPreferences(
	ctx context.Context, 
	userProfileID int64, 
	patch map[string]interface{}) error {
		
	result := r.db.GormDB.WithContext(ctx).
		Model(&entity.PartnerPreference{}).
		Where("user_profile_id = ?", userProfileID).
		Updates(patch)
	
	if result.Error != nil {
		return result.Error
	}
	
	return nil
}