package entity

import (
	"time"

	"github.com/google/uuid"
)

type UserVideo struct {
	ID              uint64    `gorm:"primaryKey;autoIncrement"`
	UserID          uuid.UUID `gorm:"type:uuid;not null;unique;index:idx_user_videos_user_id"`
	VideoURL        string    `gorm:"type:varchar(500);not null"`
	ObjectKey       string    `gorm:"type:varchar(500);not null"`
	FileName        string    `gorm:"type:varchar(255);not null"`
	FileSize        int64     `gorm:"not null"`
	DurationSeconds *int      `gorm:""` // optional
	CreatedAt       time.Time `gorm:"autoCreateTime"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime"`
}

func (UserVideo) TableName() string {
	return "user_videos"
}
