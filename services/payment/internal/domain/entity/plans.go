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
