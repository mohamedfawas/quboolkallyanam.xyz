package entity

import (
	"fmt"
	"time"

	"gorm.io/gorm"
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

func (p *PartnerPreference) ValidateAgeRange() error {
	if p.MinAgeYears != nil && p.MaxAgeYears != nil {
		if *p.MaxAgeYears < *p.MinAgeYears {
			return fmt.Errorf("max age cannot be less than min age")
		}
	}
	return nil
}

func (p *PartnerPreference) ValidateHeightRange() error {
	if p.MinHeightCm != nil && p.MaxHeightCm != nil {
		if *p.MaxHeightCm < *p.MinHeightCm {
			return fmt.Errorf("max height cannot be less than min height")
		}
	}
	return nil
}

func (p *PartnerPreference) BeforeCreate(tx *gorm.DB) error {
	if err := p.ValidateAgeRange(); err != nil {
		return err
	}
	return p.ValidateHeightRange()
}

func (p *PartnerPreference) BeforeUpdate(tx *gorm.DB) error {
	if err := p.ValidateAgeRange(); err != nil {
		return err
	}
	return p.ValidateHeightRange()
}

func (u *UserProfile) GetAge() *int {
	if u.DateOfBirth == nil {
		return nil
	}
	age := int(time.Since(*u.DateOfBirth).Hours() / 24 / 365)
	return &age
}

func (u *UserProfile) MatchesPreferences(prefs *PartnerPreference) bool {
	// Age check
	if age := u.GetAge(); age != nil {
		if prefs.MinAgeYears != nil && *age < *prefs.MinAgeYears {
			return false
		}
		if prefs.MaxAgeYears != nil && *age > *prefs.MaxAgeYears {
			return false
		}
	}

	// Height check
	if u.HeightCm != nil {
		if prefs.MinHeightCm != nil && *u.HeightCm < *prefs.MinHeightCm {
			return false
		}
		if prefs.MaxHeightCm != nil && *u.HeightCm > *prefs.MaxHeightCm {
			return false
		}
	}

	// Physical challenge check
	if u.PhysicallyChallenged && !prefs.AcceptPhysicallyChallenged {
		return false
	}

	return true
}
