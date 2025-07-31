package entity

import (
	"time"

	"github.com/google/uuid"
)

type UserProjection struct {
	UserUUID      uuid.UUID `json:"user_uuid" gorm:"type:uuid;primaryKey;column:user_uuid"`
	UserProfileID int64     `json:"user_profile_id" gorm:"uniqueIndex:idx_user_projection_user_profile_id;not null;column:user_profile_id"`
	Email         string    `json:"email" gorm:"type:varchar(255);unique;column:email"`
	FullName      string    `json:"full_name" gorm:"type:varchar(255);not null;column:full_name"`
	CreatedAt     time.Time `json:"created_at" gorm:"not null;column:created_at"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"not null;column:updated_at"`
}

func (UserProjection) TableName() string {
	return "user_projection"
}
