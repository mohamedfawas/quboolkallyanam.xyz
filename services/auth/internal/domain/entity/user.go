package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID            uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Email         string         `gorm:"size:255;not null;uniqueIndex:users_email_unique_active,where:deleted_at IS NULL"`
	Phone         string         `gorm:"size:20;not null;uniqueIndex:users_phone_unique_active,where:deleted_at IS NULL"`
	PasswordHash  string         `gorm:"size:255;not null"`
	EmailVerified bool           `gorm:"not null;default:false"`
	PremiumUntil  *time.Time     `gorm:"default:null"`
	LastLoginAt   *time.Time     `gorm:"column:last_login_at;default:null"`
	IsActive      bool           `gorm:"not null;default:true"`
	IsBlocked     bool           `gorm:"not null;default:false"`
	CreatedAt     time.Time      `gorm:"not null"`
	UpdatedAt     time.Time      `gorm:"not null"`
	DeletedAt     gorm.DeletedAt `gorm:"index;column:deleted_at"`
}

func (User) TableName() string {
	return "users"
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

func (u *User) IsPremium() bool {
	if u.PremiumUntil == nil {
		return false
	}
	return time.Now().Before(*u.PremiumUntil)
}
