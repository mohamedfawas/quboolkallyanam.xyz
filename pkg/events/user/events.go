package user

import (
	"time"

	"github.com/google/uuid"
)

type UserProfileCreatedEvent struct {
	UserID        uuid.UUID `json:"user_id"`
	UserProfileID int64     `json:"user_profile_id"`
	Email         string    `json:"email"`
	Phone         string    `json:"phone"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
