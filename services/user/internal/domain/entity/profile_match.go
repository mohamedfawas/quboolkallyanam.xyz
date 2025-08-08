package entity

import (
	"time"

	"github.com/google/uuid"
)

type ProfileMatch struct {
	ID        int64    `gorm:"primaryKey;autoIncrement"`
	UserID    uuid.UUID `gorm:"type:uuid;not null"`
	TargetID  uuid.UUID `gorm:"type:uuid;not null"`
	IsLiked   bool      `gorm:"not null;default:false"`
	IsDeleted bool      `gorm:"not null;default:false"`
	DeletedAt *time.Time
	CreatedAt time.Time `gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt time.Time `gorm:"type:timestamptz;not null;default:now()"`
}

func (ProfileMatch) TableName() string {
	return "profile_matches"
}
