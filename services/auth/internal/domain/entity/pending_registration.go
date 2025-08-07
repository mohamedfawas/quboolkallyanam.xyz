package entity

import "time"

type PendingRegistration struct {
	ID           int       `json:"id" gorm:"primaryKey"`
	Email        string    `json:"email" gorm:"uniqueIndex;size:255"`
	Phone        string    `json:"phone" gorm:"uniqueIndex;size:20"`
	PasswordHash string    `json:"-"` // Never expose password hash
	CreatedAt    time.Time `json:"created_at" gorm:"type:timestamptz;not null"`
	ExpiresAt    time.Time `json:"expires_at" gorm:"type:timestamptz;not null"`
}

func (PendingRegistration) TableName() string {
	return "pending_registrations"
}

type UserRegistrationRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Phone    string `json:"phone" validate:"required,phone"`
	Password string `json:"password" validate:"required,min=8"`
}
