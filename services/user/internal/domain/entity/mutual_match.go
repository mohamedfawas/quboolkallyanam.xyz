package entity

import (
	"time"

	"github.com/google/uuid"
)

type MutualMatch struct {
	ID        int64     `gorm:"primaryKey;autoIncrement"`
	UserID1   uuid.UUID  `gorm:"column:user_id_1;type:uuid;not null"`
	UserID2   uuid.UUID  `gorm:"column:user_id_2;type:uuid;not null"`
	MatchedAt time.Time  `gorm:"not null;default:now()"`
	IsDeleted bool       `gorm:"not null;default:false"`
	DeletedAt *time.Time `gorm:"type:timestamptz"`
	CreatedAt time.Time  `gorm:"not null;default:now()"`
	UpdatedAt time.Time  `gorm:"not null;default:now()"`
}

func (MutualMatch) TableName() string {
	return "mutual_matches"
}
