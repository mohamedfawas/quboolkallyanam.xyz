package entity

import (
	"time"

	"github.com/google/uuid"
)

type UserImage struct {
	ID           uint      `gorm:"primaryKey;autoIncrement"`
	UserID       uuid.UUID `gorm:"type:uuid;not null;index:idx_user_images_user_id"`
	ImageURL     string    `gorm:"type:varchar(500);not null"`
	ObjectKey    string    `gorm:"type:varchar(500);not null"`
	DisplayOrder int16     `gorm:"type:smallint;not null;check:display_order >= 1 AND display_order <= 3"`
	CreatedAt    time.Time `gorm:"not null;default:now()"`
	UpdatedAt    time.Time `gorm:"not null;default:now()"`
}

func (UserImage) TableName() string {
	return "user_images"
}
