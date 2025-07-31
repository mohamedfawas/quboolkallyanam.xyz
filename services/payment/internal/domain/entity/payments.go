package entity

import (
	"time"
)

type Payment struct {
	ID                int64     `gorm:"primaryKey;autoIncrement"`
	UserID            string    `gorm:"type:uuid;not null"`
	PlanID            string    `gorm:"type:varchar(100);not null"`
	RazorpayOrderID   string    `gorm:"type:varchar(255);not null"`
	RazorpayPaymentID string    `gorm:"type:varchar(255)"`
	RazorpaySignature string    `gorm:"type:varchar(255)"`
	Amount            float64   `gorm:"type:decimal(10,2);not null"`
	Currency          string    `gorm:"type:varchar(3);not null"`
	Status            string    `gorm:"type:varchar(50);not null;default:'pending'"`
	PaymentMethod     string    `gorm:"type:varchar(50);not null;default:'razorpay'"`
	ExpiresAt         time.Time `gorm:"not null"`
	CreatedAt         time.Time `gorm:"default:current_timestamp"`
	UpdatedAt         time.Time `gorm:"default:current_timestamp"`
}

func (Payment) TableName() string {
	return "payments"
}

type PaymentOrderResponse struct {
	RazorpayOrderID string
	Amount          float64
	Currency        string
	PlanID          string
	ExpiresAt       time.Time
}

type VerifyPaymentRequest struct {
	RazorpayOrderID   string `json:"razorpay_order_id" binding:"required"`
	RazorpayPaymentID string `json:"razorpay_payment_id" binding:"required"`
	RazorpaySignature string `json:"razorpay_signature" binding:"required"`
}

type VerifyPaymentResponse struct {
	SubscriptionID        string    `json:"subscription_id"`
	SubscriptionStartDate time.Time `json:"subscription_start_date"`
	SubscriptionEndDate   time.Time `json:"subscription_end_date"`
	SubscriptionStatus    string    `json:"subscription_status"`
}

type ShowPaymentPageResponse struct {
	RazorpayOrderID    string `json:"razorpay_order_id"`
	RazorpayKeyID      string `json:"razorpay_key_id"`
	PlanID             string `json:"plan_id"`
	Amount             int64  `json:"amount"` // Amount in paise
	DisplayAmount      string `json:"display_amount"`
	PlanDurationInDays int32  `json:"plan_duration_in_days"`
}

type GetPaymentHistoryResponse struct {
	ID              int64     `json:"id"`
	PlanID          string    `json:"plan_id"`
	RazorpayOrderID string    `json:"razorpay_order_id"`
	Amount          float64   `json:"amount"`
	Currency        string    `json:"currency"`
	Status          string    `json:"status"`
	PaymentMethod   string    `json:"payment_method"`
	CreatedAt       time.Time `json:"created_at"`
}
