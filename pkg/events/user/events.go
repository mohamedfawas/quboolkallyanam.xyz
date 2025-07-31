package user

import (
	"time"

	"github.com/google/uuid"
)

type UserProfileUpdatedEvent struct {
	UserID        uuid.UUID `json:"user_id"`
	UserProfileID int64     `json:"user_profile_id"`
	FullName      string    `json:"full_name"`
	Email         string    `json:"email"`
	Phone         string    `json:"phone"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
