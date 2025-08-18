package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/validation"
	"gorm.io/gorm"
)

type UserProfile struct {
	ID                   int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID               uuid.UUID `json:"user_id" gorm:"type:uuid;not null;uniqueIndex"`
	IsBride              bool      `json:"is_bride" gorm:"not null;default:false"`
	FullName             string    `json:"full_name" gorm:"size:200"`
	Email                string    `json:"email" gorm:"not null;size:255"`
	Phone                string    `json:"phone" gorm:"not null;size:20"`
	DateOfBirth          time.Time `json:"date_of_birth" gorm:"type:date"`
	HeightCm             int16     `json:"height_cm" gorm:"type:smallint;not null;check:height_cm > 0"`
	PhysicallyChallenged bool      `json:"physically_challenged" gorm:"not null;default:false"`
	ProfileCompleted     bool      `json:"profile_completed" gorm:"not null;default:false"`

	Community             validation.Community      `json:"community" gorm:"type:varchar(255)"`
	MaritalStatus         validation.MaritalStatus  `json:"marital_status" gorm:"type:varchar(255)"`
	Profession            validation.Profession     `json:"profession" gorm:"type:varchar(255)"`
	ProfessionType        validation.ProfessionType `json:"profession_type" gorm:"type:varchar(255)"`
	HighestEducationLevel validation.EducationLevel `json:"highest_education_level" gorm:"type:varchar(255)"`
	HomeDistrict          validation.HomeDistrict   `json:"home_district" gorm:"type:varchar(255)"`

	ProfileImageKey string         `json:"profile_image_key" gorm:"size:255"`
	LastLogin       time.Time      `json:"last_login" gorm:"not null;default:CURRENT_TIMESTAMP"`
	CreatedAt       time.Time      `json:"created_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt       time.Time      `json:"updated_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
	DeletedAt       gorm.DeletedAt `gorm:"index"            json:"deleted_at,omitempty"`

	// Association
	PartnerPreference PartnerPreference `json:"partner_preference,omitempty" gorm:"foreignKey:UserProfileID"`
}

func (UserProfile) TableName() string {
	return "user_profiles"
}

// This should be similar to the proto file code.
type UpdateUserProfileRequest struct {
	IsBride               *bool   `json:"is_bride"`
	FullName              *string `json:"full_name"`
	DateOfBirth           *string `json:"date_of_birth"`
	HeightCm              *int32  `json:"height_cm"`
	PhysicallyChallenged  *bool   `json:"physically_challenged"`
	Community             *string `json:"community"`
	MaritalStatus         *string `json:"marital_status"`
	Profession            *string `json:"profession"`
	ProfessionType        *string `json:"profession_type"`
	HighestEducationLevel *string `json:"highest_education_level"`
	HomeDistrict          *string `json:"home_district"`
	ProfileCompleted      *bool   `json:"profile_completed"`
}

type UserProfileResponse struct {
	ID                int64   `json:"id"`
	FullName          string  `json:"full_name"`
	ProfilePictureURL *string `json:"profile_picture_url"`
	Age               int32   `json:"age"`
	HeightCm          int32   `json:"height_cm"`
	MaritalStatus     string  `json:"marital_status"`
	Profession        string  `json:"profession"`
	HomeDistrict      string  `json:"home_district"`
}
