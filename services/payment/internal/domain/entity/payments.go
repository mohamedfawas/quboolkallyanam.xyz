package entity

import (
	"time"
)

type Payment struct {
	ID                int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	OrderID           string    `gorm:"type:varchar(255);uniqueIndex:idx_payments_order_id;not null" json:"order_id"`
	UserID            string    `gorm:"type:uuid;not null;index:idx_payments_user_id" json:"user_id"`
	PlanID            string    `gorm:"type:varchar(100);not null" json:"plan_id"`
	RazorpayOrderID   string    `gorm:"type:varchar(255);uniqueIndex:idx_payments_razorpay_order_id;not null" json:"razorpay_order_id"`
	RazorpayPaymentID string    `gorm:"type:varchar(255)" json:"razorpay_payment_id,omitempty"`
	RazorpaySignature string    `gorm:"type:varchar(255)" json:"razorpay_signature,omitempty"`
	Amount            float64   `gorm:"type:decimal(10,2);not null;check:amount > 0" json:"amount"`
	Currency          string    `gorm:"type:varchar(3);not null" json:"currency"`
	Status            string    `gorm:"type:varchar(50);not null;default:'pending';check:status IN ('pending','completed','failed');index:idx_payments_status" json:"status"`
	PaymentMethod     string    `gorm:"type:varchar(50);not null;default:'razorpay';check:payment_method IN ('razorpay')" json:"payment_method"`
	ExpiresAt         time.Time `gorm:"not null" json:"expires_at"`
	CreatedAt         time.Time `gorm:"default:current_timestamp" json:"created_at"`
	UpdatedAt         time.Time `gorm:"default:current_timestamp" json:"updated_at"`
}

func (Payment) TableName() string {
	return "payments"
}

type PaymentOrderResponse struct {
	OrderID         string
	RazorpayOrderID string
	RazorpayKeyID   string
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
