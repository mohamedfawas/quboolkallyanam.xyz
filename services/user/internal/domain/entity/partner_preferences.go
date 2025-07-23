package entity

import (
	"time"
)

type PartnerPreference struct {
	ID            int64 `json:"id" gorm:"primaryKey"`
	UserProfileID int64 `json:"user_profile_id" gorm:"not null;uniqueIndex;constraint:OnDelete:CASCADE"`

	MinAgeYears *int `json:"min_age_years" gorm:"check:min_age_years BETWEEN 18 AND 100"`
	MaxAgeYears *int `json:"max_age_years" gorm:"check:max_age_years BETWEEN 18 AND 100"`
	MinHeightCm *int `json:"min_height_cm" gorm:"check:min_height_cm BETWEEN 130 AND 220"`
	MaxHeightCm *int `json:"max_height_cm" gorm:"check:max_height_cm BETWEEN 130 AND 220"`

	AcceptPhysicallyChallenged bool `json:"accept_physically_challenged" gorm:"not null;default:true"`

	PreferredCommunities     []CommunityEnum      `json:"preferred_communities" gorm:"type:jsonb;serializer:json;default:'[]'"`
	PreferredMaritalStatus   []MaritalStatusEnum  `json:"preferred_marital_status" gorm:"type:jsonb;serializer:json;default:'[]'"`
	PreferredProfessions     []ProfessionEnum     `json:"preferred_professions" gorm:"type:jsonb;serializer:json;default:'[]'"`
	PreferredProfessionTypes []ProfessionTypeEnum `json:"preferred_profession_types" gorm:"type:jsonb;serializer:json;default:'[]'"`
	PreferredEducationLevels []EducationLevelEnum `json:"preferred_education_levels" gorm:"type:jsonb;serializer:json;default:'[]'"`
	PreferredHomeDistricts   []HomeDistrictEnum   `json:"preferred_home_districts" gorm:"type:jsonb;serializer:json;default:'[]'"`

	CreatedAt time.Time `json:"created_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `json:"updated_at" gorm:"not null;default:CURRENT_TIMESTAMP"`

	// Association
	UserProfile *UserProfile `json:"user_profile,omitempty" gorm:"foreignKey:UserProfileID"`
}

func (PartnerPreference) TableName() string {
	return "partner_preferences"
}

type UpdateUserPartnerPreferencesRequest struct {
	MinAgeYears                *int      `json:"min_age_years,omitempty"`
	MaxAgeYears                *int      `json:"max_age_years,omitempty"`
	MinHeightCM                *int      `json:"min_height_cm,omitempty"`
	MaxHeightCM                *int      `json:"max_height_cm,omitempty"`
	AcceptPhysicallyChallenged *bool     `json:"accept_physically_challenged,omitempty"`
	PreferredCommunities       *[]string `json:"preferred_communities,omitempty"`
	PreferredMaritalStatus     *[]string `json:"preferred_marital_status,omitempty"`
	PreferredProfessions       *[]string `json:"preferred_professions,omitempty"`
	PreferredProfessionTypes   *[]string `json:"preferred_profession_types,omitempty"`
	PreferredEducationLevels   *[]string `json:"preferred_education_levels,omitempty"`
	PreferredHomeDistricts     *[]string `json:"preferred_home_districts,omitempty"`
}
