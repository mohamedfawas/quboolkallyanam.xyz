package entity

import (
	"time"
)

type SubscriptionPlan struct {
	ID           string    `gorm:"type:varchar(100);primaryKey" json:"id"`
	DurationDays int       `gorm:"not null" json:"duration_days"`
	Amount       float64   `gorm:"type:decimal(10,2);not null" json:"amount"`
	Currency     string    `gorm:"type:varchar(3);not null" json:"currency"`
	Description  string    `gorm:"type:text" json:"description"`
	IsActive     bool      `gorm:"not null;default:true" json:"is_active"`
	CreatedAt    time.Time `gorm:"not null;default:current_timestamp" json:"created_at"`
	UpdatedAt    time.Time `gorm:"not null;default:current_timestamp" json:"updated_at"`
}

func (SubscriptionPlan) TableName() string {
	return "subscription_plans"
}

type CreateSubscriptionPlanRequest struct {
	ID           string  `json:"id"`
	DurationDays int     `json:"duration_days"`
	Amount       float64 `json:"amount"`
	Currency     string  `json:"currency"`
	Description  string  `json:"description"`
	IsActive     bool    `json:"is_active"`
}

type UpdateSubscriptionPlanRequest struct {
	ID           string   `json:"id"`
	DurationDays *int     `json:"duration_days,omitempty"`
	Amount       *float64 `json:"amount,omitempty"`
	Currency     *string  `json:"currency,omitempty"`
	Description  *string  `json:"description,omitempty"`
	IsActive     *bool    `json:"is_active,omitempty"`
}
