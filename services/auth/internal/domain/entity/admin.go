package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Admin struct {
	ID           uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Email        string         `gorm:"size:255;not null;uniqueIndex:idx_admins_email"`
	PasswordHash string         `gorm:"size:255;not null"`
	IsActive     bool           `gorm:"not null;default:true"`
	CreatedAt    time.Time      `gorm:"not null"`
	UpdatedAt    time.Time      `gorm:"not null"`
	DeletedAt    gorm.DeletedAt `gorm:"index;column:deleted_at"`
}

func (Admin) TableName() string {
	return "admins"
}

func (a *Admin) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}
