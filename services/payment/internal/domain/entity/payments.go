package entity

import (
	"time"
)

type Payment struct {
	ID                int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID            string    `gorm:"type:uuid;not null;index:idx_payments_user_id" json:"user_id"`
	RazorpayOrderID   string    `gorm:"type:varchar(255);uniqueIndex:idx_payments_razorpay_order_id;not null" json:"razorpay_order_id"`
	RazorpayPaymentID string    `gorm:"type:varchar(255)" json:"razorpay_payment_id,omitempty"`
	RazorpaySignature string    `gorm:"type:varchar(255)" json:"razorpay_signature,omitempty"`
	Amount            float64   `gorm:"type:decimal(10,2);not null;check:amount > 0" json:"amount"`
	Currency          string    `gorm:"type:varchar(3);not null" json:"currency"`
	Status            string    `gorm:"type:varchar(50);not null;default:'pending';check:status IN ('pending','completed','failed');index:idx_payments_status" json:"status"`
	PaymentMethod     string    `gorm:"type:varchar(50);not null;default:'razorpay';check:payment_method IN ('razorpay')" json:"payment_method"`
	CreatedAt         time.Time `gorm:"default:current_timestamp" json:"created_at"`
	UpdatedAt         time.Time `gorm:"default:current_timestamp" json:"updated_at"`
}

func (Payment) TableName() string {
	return "payments"
}
