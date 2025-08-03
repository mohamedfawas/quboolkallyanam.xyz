package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/database/postgres"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/validation"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/user/internal/domain/entity"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/user/internal/domain/repository"
	"gorm.io/gorm"
)

type userProfileRepository struct {
	db *postgres.Client
}

func NewUserProfileRepository(db *postgres.Client) repository.UserProfileRepository {
	return &userProfileRepository{db: db}
}

func (r *userProfileRepository) CreateUserProfile(ctx context.Context, userProfile *entity.UserProfile) error {
	return r.db.GormDB.WithContext(ctx).
		Create(userProfile).Error
}

func (r *userProfileRepository) UpdateLastLogin(ctx context.Context, userID uuid.UUID) error {
	now := time.Now().UTC()
	return r.db.GormDB.WithContext(ctx).
		Model(&entity.UserProfile{}).
		Where("user_id = ?", userID).
		Updates(map[string]interface{}{
			"last_login": now,
			"updated_at": now,
		}).Error
}

func (r *userProfileRepository) ProfileExists(ctx context.Context, userID uuid.UUID) (bool, error) {
	var matched int64
	err := r.db.GormDB.WithContext(ctx).
		Model(&entity.UserProfile{}).
		Where("user_id = ?", userID).
		Count(&matched).Error
	if err != nil {
		return false, err
	}
	return matched > 0, nil
}

func (r *userProfileRepository) GetProfileByUserID(ctx context.Context, userID uuid.UUID) (*entity.UserProfile, error) {
	var profile entity.UserProfile
	err := r.db.GormDB.WithContext(ctx).
		Model(&entity.UserProfile{}).
		Where("user_id = ?", userID).
		First(&profile).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &profile, nil
}

func (r *userProfileRepository) UpdateUserProfile(ctx context.Context, userProfile *entity.UserProfile) error {
	return r.db.GormDB.WithContext(ctx).
		Model(&entity.UserProfile{}).
		Where("user_id = ?", userProfile.UserID).
		Updates(userProfile).Error
}

func (r *userProfileRepository) GetUserProfileByID(ctx context.Context, id uint) (*entity.UserProfile, error) {
	var profile entity.UserProfile
	err := r.db.GormDB.WithContext(ctx).
		Model(&entity.UserProfile{}).
		Where("id = ?", id).
		First(&profile).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &profile, nil
}

func (r *userProfileRepository) GetPotentialProfiles(
	ctx context.Context,
	userID uuid.UUID,
	excludedIDs []uuid.UUID,
	preferences *entity.PartnerPreference,
	limit int,
	offset int,
	isUserBride bool) ([]*entity.UserProfile, int64, error) {

	// 1) Build base query
    baseQuery := r.db.GormDB.WithContext(ctx).
        Model(&entity.UserProfile{}).
        Where("user_id <> ?", userID).
        Where("is_bride = ?", !isUserBride)

	if len(excludedIDs) > 0 {
		baseQuery = baseQuery.Where("user_id NOT IN (?)", excludedIDs)
	}

	// Apply preferences filters if provided
	if preferences != nil {
        baseQuery = r.applyPreferencesFilters(baseQuery, preferences)
    }

	// 2) Count total matches
    var total int64
    if err := baseQuery.
        Session(&gorm.Session{}). // clone so we don’t carry over Limit/Offset/etc.
        Count(&total).Error; err != nil {
        return nil, 0, err
    }

	// 3) Fetch paginated results
    var profiles []*entity.UserProfile
    if err := baseQuery.
        Order("last_login DESC").
        Limit(limit).
        Offset(offset).
        Find(&profiles).Error; err != nil {
        return nil, 0, err
    }

	return profiles, total, nil
}

func (r *userProfileRepository) applyPreferencesFilters(
	query *gorm.DB,
	preferences *entity.PartnerPreference) *gorm.DB {

	// age filtering
	now := time.Now().UTC()
	maxBirthDate := now.AddDate(-int(preferences.MinAgeYears), 0, 0)
	minBirthDate := now.AddDate(-int(preferences.MaxAgeYears)-1, 0, 0)
	query = query.Where("date_of_birth BETWEEN ? AND ?", minBirthDate, maxBirthDate)

	// Height filtering
	query = query.Where("height_cm BETWEEN ? AND ?",
		preferences.MinHeightCm, preferences.MaxHeightCm)

	// Physically challenged filtering
	if !preferences.AcceptPhysicallyChallenged {
		query = query.Where("physically_challenged = ?", false)
	}

	// Community filtering, if "any" value is there, it will be the only value present there
	if preferences.PreferredCommunities[0] != validation.CommunityAny {
		query = query.Where("community IN ?", preferences.PreferredCommunities)
	}

	if preferences.PreferredMaritalStatus[0] != validation.MaritalStatusAny {
		query = query.Where("marital_status IN ?", preferences.PreferredMaritalStatus)
	}

	if preferences.PreferredProfessions[0] != validation.ProfessionAny {
		query = query.Where("profession IN ?", preferences.PreferredProfessions)
	}

	if preferences.PreferredProfessionTypes[0] != validation.ProfessionTypeAny {
		query = query.Where("profession_type IN ?", preferences.PreferredProfessionTypes)
	}

	if preferences.PreferredEducationLevels[0] != validation.EducationLevelAny {
		query = query.Where("highest_education_level IN ?", preferences.PreferredEducationLevels)
	}

	if preferences.PreferredHomeDistricts[0] != validation.HomeDistrictAny {
		query = query.Where("home_district IN ?", preferences.PreferredHomeDistricts)
	}

	return query
}
