package entity

import (
    "time"

    "github.com/google/uuid"
    "gorm.io/gorm"
)

type UserImage struct {
    ID           int64           `gorm:"primaryKey;autoIncrement"`
    UserID       uuid.UUID      `gorm:"type:uuid;not null;index:idx_user_images_user_id"`
    ObjectKey    string         `gorm:"type:varchar(500);not null"`
    DisplayOrder int16          `gorm:"type:smallint;not null;check:display_order >= 1 AND display_order <= 3;index:uniq_active_img_order,unique"`
    CreatedAt    time.Time      `gorm:"not null;default:now()"`
    UpdatedAt    time.Time      `gorm:"not null;default:now()"`
    DeletedAt    gorm.DeletedAt `gorm:"index"`  
}

func (UserImage) TableName() string {
    return "user_images"
}
