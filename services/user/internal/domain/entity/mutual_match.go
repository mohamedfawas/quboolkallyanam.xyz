package entity

import (
	"time"

	"github.com/google/uuid"
)

type MutualMatch struct {
	ID        uint64     `gorm:"primaryKey;autoIncrement"`
	UserID1   uuid.UUID  `gorm:"type:uuid;not null"`
	UserID2   uuid.UUID  `gorm:"type:uuid;not null"`
	MatchedAt time.Time  `gorm:"not null;default:now()"`
	IsActive  bool       `gorm:"not null;default:true"`
	IsDeleted bool       `gorm:"not null;default:false"`
	DeletedAt *time.Time `gorm:"type:timestamptz"`
	CreatedAt time.Time  `gorm:"not null;default:now()"`
	UpdatedAt time.Time  `gorm:"not null;default:now()"`
}

func (MutualMatch) TableName() string {
	return "mutual_matches"
}
