package authevents

import (
	"github.com/google/uuid"
)

type UserOTPRequestedEvent struct {
	Email         string `json:"email"`
	OTP           string `json:"otp"`
	ExpiryMinutes int    `json:"expiry_minutes"`
}

type UserLoginSuccessEvent struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	Phone  string    `json:"phone"`
}
