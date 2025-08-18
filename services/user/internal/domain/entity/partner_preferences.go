package entity

import (
	"time"

	"github.com/lib/pq"
)

type PartnerPreference struct {
	ID            int64 `json:"id" gorm:"primaryKey"`
	UserProfileID int64 `json:"user_profile_id" gorm:"not null;uniqueIndex"`

	// uint16 better maps to smallint than uint8
	MinAgeYears int16 `json:"min_age_years" gorm:"type:smallint;not null;default:18"`
	MaxAgeYears int16 `json:"max_age_years" gorm:"type:smallint;not null;default:100"`
	MinHeightCm int16 `json:"min_height_cm" gorm:"type:smallint;not null;default:130"`
	MaxHeightCm int16 `json:"max_height_cm" gorm:"type:smallint;not null;default:220"`

	AcceptPhysicallyChallenged bool `json:"accept_physically_challenged" gorm:"not null;default:true"`

	PreferredCommunities     pq.StringArray `json:"preferred_communities" gorm:"type:text[];not null;default:'{}'"`
	PreferredMaritalStatus   pq.StringArray `json:"preferred_marital_status" gorm:"type:text[];not null;default:'{}'"`
	PreferredProfessions     pq.StringArray `json:"preferred_professions" gorm:"type:text[];not null;default:'{}'"`
	PreferredProfessionTypes pq.StringArray `json:"preferred_profession_types" gorm:"type:text[];not null;default:'{}'"`
	PreferredEducationLevels pq.StringArray `json:"preferred_education_levels" gorm:"type:text[];not null;default:'{}'"`
	PreferredHomeDistricts   pq.StringArray `json:"preferred_home_districts" gorm:"type:text[];not null;default:'{}'"`

	CreatedAt time.Time `json:"created_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `json:"updated_at" gorm:"not null;default:CURRENT_TIMESTAMP"`

	// Association
	UserProfile *UserProfile `json:"user_profile,omitempty" gorm:"foreignKey:UserProfileID"`
}

func (PartnerPreference) TableName() string {
	return "partner_preferences"
}



	
type UpdateUserPartnerPreferencesRequest struct {
	MinAgeYears                *int16   `json:"min_age_years,omitempty"`
	MaxAgeYears                *int16   `json:"max_age_years,omitempty"`
	MinHeightCM                *int16   `json:"min_height_cm,omitempty"`
	MaxHeightCM                *int16   `json:"max_height_cm,omitempty"`
	AcceptPhysicallyChallenged *bool     `json:"accept_physically_challenged,omitempty"`
	PreferredCommunities       *[]string `json:"preferred_communities,omitempty"`
	PreferredMaritalStatus     *[]string `json:"preferred_marital_status,omitempty"`
	PreferredProfessions       *[]string `json:"preferred_professions,omitempty"`
	PreferredProfessionTypes   *[]string `json:"preferred_profession_types,omitempty"`
	PreferredEducationLevels   *[]string `json:"preferred_education_levels,omitempty"`
	PreferredHomeDistricts     *[]string `json:"preferred_home_districts,omitempty"`
	AcceptAllCommunities       *bool     `json:"accept_all_communities,omitempty"`
	AcceptAllMaritalStatus     *bool     `json:"accept_all_marital_status,omitempty"`
	AcceptAllProfessions       *bool     `json:"accept_all_professions,omitempty"`
	AcceptAllProfessionTypes   *bool     `json:"accept_all_profession_types,omitempty"`
	AcceptAllEducationLevels   *bool     `json:"accept_all_education_levels,omitempty"`
	AcceptAllHomeDistricts     *bool     `json:"accept_all_home_districts,omitempty"`
}
