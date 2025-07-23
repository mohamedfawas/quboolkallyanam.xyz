package entity

import (
	"time"

	"github.com/google/uuid"
)

type UserProfile struct {
	ID                    int64               `json:"id" gorm:"primaryKey"`
	UserID                uuid.UUID           `json:"user_id" gorm:"type:uuid;not null;uniqueIndex"`
	IsBride               bool                `json:"is_bride" gorm:"not null;default:false"`
	FullName              *string             `json:"full_name" gorm:"size:200"`
	Email                 *string             `json:"email" gorm:"size:255"`
	Phone                 *string             `json:"phone" gorm:"size:20"`
	DateOfBirth           *time.Time          `json:"date_of_birth" gorm:"type:date"`
	HeightCm              *int                `json:"height_cm" gorm:"check:height_cm BETWEEN 130 AND 220"`
	PhysicallyChallenged  bool                `json:"physically_challenged" gorm:"not null;default:false"`
	Community             *CommunityEnum      `json:"community" gorm:"type:community_enum"`
	MaritalStatus         *MaritalStatusEnum  `json:"marital_status" gorm:"type:marital_status_enum"`
	Profession            *ProfessionEnum     `json:"profession" gorm:"type:profession_enum"`
	ProfessionType        *ProfessionTypeEnum `json:"profession_type" gorm:"type:profession_type_enum"`
	HighestEducationLevel *EducationLevelEnum `json:"highest_education_level" gorm:"type:education_level_enum"`
	HomeDistrict          *HomeDistrictEnum   `json:"home_district" gorm:"type:home_district_enum"`
	ProfilePictureURL     *string             `json:"profile_picture_url" gorm:"size:255"`
	LastLogin             time.Time           `json:"last_login" gorm:"not null;default:CURRENT_TIMESTAMP"`
	CreatedAt             time.Time           `json:"created_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt             time.Time           `json:"updated_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
	IsDeleted             bool                `json:"is_deleted" gorm:"not null;default:false"`
	DeletedAt             *time.Time          `json:"deleted_at"`

	// Association
	PartnerPreference *PartnerPreference `json:"partner_preference,omitempty" gorm:"foreignKey:UserProfileID"`
}

func (UserProfile) TableName() string {
	return "user_profiles"
}

type UpdateUserProfileRequest struct {
	IsBride               *bool      `json:"is_bride"`
	FullName              *string    `json:"full_name"`
	DateOfBirth           *time.Time `json:"date_of_birth"`
	HeightCm              *int       `json:"height_cm"`
	PhysicallyChallenged  *bool      `json:"physically_challenged"`
	Community             *string    `json:"community"`
	MaritalStatus         *string    `json:"marital_status"`
	Profession            *string    `json:"profession"`
	ProfessionType        *string    `json:"profession_type"`
	HighestEducationLevel *string    `json:"highest_education_level"`
	HomeDistrict          *string    `json:"home_district"`
}
