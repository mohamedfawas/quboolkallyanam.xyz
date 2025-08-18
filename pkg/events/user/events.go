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

type UserInterestSentEvent struct {
	ReceiverEmail   string `json:"receiver_user_id"`
	SenderProfileID int64  `json:"sender_profile_id"`
	SenderName      string `json:"sender_name"`
}

type MutualMatchCreatedEvent struct {
	User1Email     string `json:"user1_email"`
	User1ProfileID int64  `json:"user1_profile_id"`
	User1FullName  string `json:"user1_full_name"`
	User2Email     string `json:"user2_email"`
	User2ProfileID int64  `json:"user2_profile_id"`
	User2FullName  string `json:"user2_full_name"`
}
