package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/database/postgres"
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

func (r *userProfileRepository) UpdateLastLogin(ctx context.Context, userID uuid.UUID, lastLogin time.Time) error {
	return r.db.GormDB.WithContext(ctx).
		Model(&entity.UserProfile{}).
		Where("user_id = ?", userID).
		Updates(map[string]interface{}{
			"last_login": lastLogin,
			"updated_at": time.Now().UTC(),
		}).Error
}

func (r *userProfileRepository) ProfileExists(ctx context.Context, userID uuid.UUID) (bool, error) {
	var matched int64
	err := r.db.GormDB.WithContext(ctx).
		Model(&entity.UserProfile{}).
		Where("user_id = ? AND is_deleted = false", userID).
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
		Where("user_id = ? AND is_deleted = false", userID).
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
		Where("id = ? AND is_deleted = false", id).
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
	isUserBride bool) ([]*entity.UserProfile, error) {

	query := r.db.GormDB.WithContext(ctx).
		Model(&entity.UserProfile{}).
		Where("user_id != ? AND is_deleted = ?", userID, false)

	if len(excludedIDs) > 0 {
		query = query.Where("user_id NOT IN ?", excludedIDs)
	}

	query = query.Where("is_bride = ?", !isUserBride)

	// Apply preferences filters if provided
	if preferences != nil {
		query = r.applyPreferencesFilters(query, preferences)
	}

	var profiles []*entity.UserProfile
	err := query.
		Order("last_login DESC").
		Limit(50).
		Find(&profiles).Error

	if err != nil {
		return nil, err
	}

	return profiles, nil
}

func (r *userProfileRepository) applyPreferencesFilters(
	query *gorm.DB,
	preferences *entity.PartnerPreference) *gorm.DB {

	if preferences.MinAgeYears != nil && preferences.MaxAgeYears != nil {
		now := time.Now().UTC()
		maxBirthDate := now.AddDate(-*preferences.MinAgeYears, 0, 0)
		minBirthDate := now.AddDate(-*preferences.MaxAgeYears-1, 0, 0)
		query = query.Where("date_of_birth BETWEEN ? AND ?", minBirthDate, maxBirthDate)
	}

	if preferences.MinHeightCm != nil && preferences.MaxHeightCm != nil {
		query = query.Where("height_cm BETWEEN ? AND ?",
			*preferences.MinHeightCm, *preferences.MaxHeightCm)
	}

	if !preferences.AcceptPhysicallyChallenged {
		query = query.Where("physically_challenged = ?", false)
	}

	// Community filter
	if len(preferences.PreferredCommunities) > 0 {
		communities := r.convertCommunitiesToStrings(preferences.PreferredCommunities)
		query = query.Where("community IN ?", communities)
	}

	// Marital status filter
	if len(preferences.PreferredMaritalStatus) > 0 {
		statuses := r.convertMaritalStatusToStrings(preferences.PreferredMaritalStatus)
		query = query.Where("marital_status IN ?", statuses)
	}

	// Profession filter
	if len(preferences.PreferredProfessions) > 0 {
		professions := r.convertProfessionsToStrings(preferences.PreferredProfessions)
		query = query.Where("profession IN ?", professions)
	}

	// Education level filter
	if len(preferences.PreferredEducationLevels) > 0 {
		levels := r.convertEducationLevelsToStrings(preferences.PreferredEducationLevels)
		query = query.Where("highest_education_level IN ?", levels)
	}

	// Home district filter
	if len(preferences.PreferredHomeDistricts) > 0 {
		districts := r.convertHomeDistrictsToStrings(preferences.PreferredHomeDistricts)
		query = query.Where("home_district IN ?", districts)
	}

	return query
}

// Helper methods to convert enums to strings efficiently
func (r *userProfileRepository) convertCommunitiesToStrings(communities []entity.CommunityEnum) []string {
	if len(communities) == 0 {
		return nil
	}
	result := make([]string, len(communities))
	for i, c := range communities {
		result[i] = string(c)
	}
	return result
}

func (r *userProfileRepository) convertMaritalStatusToStrings(statuses []entity.MaritalStatusEnum) []string {
	if len(statuses) == 0 {
		return nil
	}
	result := make([]string, len(statuses))
	for i, s := range statuses {
		result[i] = string(s)
	}
	return result
}

func (r *userProfileRepository) convertProfessionsToStrings(professions []entity.ProfessionEnum) []string {
	if len(professions) == 0 {
		return nil
	}
	result := make([]string, len(professions))
	for i, p := range professions {
		result[i] = string(p)
	}
	return result
}

func (r *userProfileRepository) convertEducationLevelsToStrings(levels []entity.EducationLevelEnum) []string {
	if len(levels) == 0 {
		return nil
	}
	result := make([]string, len(levels))
	for i, l := range levels {
		result[i] = string(l)
	}
	return result
}

func (r *userProfileRepository) convertHomeDistrictsToStrings(districts []entity.HomeDistrictEnum) []string {
	if len(districts) == 0 {
		return nil
	}
	result := make([]string, len(districts))
	for i, d := range districts {
		result[i] = string(d)
	}
	return result
}
