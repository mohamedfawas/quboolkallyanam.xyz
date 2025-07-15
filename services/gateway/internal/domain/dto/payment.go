package dto

import "time"

type PaymentOrderRequest struct {
	PlanID string `json:"plan_id" binding:"required"`
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

type CreatePaymentOrderResponse struct {
	OrderID         string    `json:"order_id"`
	PaymentURL      string    `json:"payment_url"`
	AmountInINR     string    `json:"amount_in_inr"`
	RazorpayOrderID string    `json:"razorpay_order_id"`
	RazorpayKeyID   string    `json:"razorpay_key_id"`
	PlanID          string    `json:"plan_id"`
	ExpiresAt       time.Time `json:"expires_at"`
}

type ShowPaymentPageRequest struct {
	RazorpayOrderID string `json:"razorpay_order_id" binding:"required"`
}

type ShowPaymentPageResponse struct {
	PlanID             string `json:"plan_id"`
	DisplayAmount      string `json:"display_amount"`
	PlanDurationInDays int    `json:"plan_duration_in_days"`
}
