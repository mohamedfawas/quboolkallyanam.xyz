package entity

import "time"

type Subscription struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    string    `gorm:"type:uuid;not null;index:idx_subscriptions_user_id" json:"user_id"`
	PlanID    string    `gorm:"type:varchar(100);not null" json:"plan_id"`
	StartDate time.Time `gorm:"not null" json:"start_date"`
	EndDate   time.Time `gorm:"not null" json:"end_date"`
	Status    string    `gorm:"type:varchar(50);not null;default:'active';check:status IN ('active','expired','cancelled');index:idx_subscriptions_status" json:"status"`
	CreatedAt time.Time `gorm:"default:current_timestamp" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:current_timestamp" json:"updated_at"`
}

func (Subscription) TableName() string {
	return "subscriptions"
}
