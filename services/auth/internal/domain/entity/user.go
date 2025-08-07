package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID            uuid.UUID      `gorm:"type:uuid;primaryKey"`
	Email         string         `gorm:"size:255;not null;uniqueIndex:users_email_unique_active,where:deleted_at IS NULL"`
	Phone         string         `gorm:"size:20;not null;uniqueIndex:users_phone_unique_active,where:deleted_at IS NULL"`
	PasswordHash  string         `gorm:"size:255;not null"`
	EmailVerified bool           `gorm:"not null;default:false"`
	PremiumUntil  *time.Time     `gorm:"type:timestamptz"`
	LastLoginAt   *time.Time     `gorm:"type:timestamptz;column:last_login_at"`
	IsBlocked     bool           `gorm:"not null;default:false;index:users_blocked_email_idx,where:is_blocked = true;index:users_blocked_phone_idx,where:is_blocked = true"`
	CreatedAt     time.Time      `gorm:"type:timestamptz;not null"`
	UpdatedAt     time.Time      `gorm:"type:timestamptz;not null"`
	DeletedAt     gorm.DeletedAt `gorm:"type:timestamptz;index;column:deleted_at"`
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


type GetUserResponse struct {
	ID            uuid.UUID      `json:"id"`
	Email         string         `json:"email"`
	Phone         string         `json:"phone"`
	EmailVerified bool           `json:"email_verified"`
	PremiumUntil  *time.Time     `json:"premium_until"`
	LastLoginAt   *time.Time     `json:"last_login_at"`
	IsBlocked     bool           `json:"is_blocked"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}