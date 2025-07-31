package dto

import "time"

type PaymentOrderRequest struct {
	PlanID string `json:"plan_id" binding:"required"`
}

type PaymentOrderResponse struct {
	RazorpayOrderID string
	Amount          float64
	Currency        string
	PlanID          string
	ExpiresAt       time.Time
}

type CreatePaymentOrderResponse struct {
	PaymentURL      string    `json:"payment_url"`
	RazorpayOrderID string    `json:"razorpay_order_id"`
	Amount          string    `json:"amount"` // Display amount like "999.00"
	PlanID          string    `json:"plan_id"`
	ExpiresAt       time.Time `json:"expires_at"`
}

type ShowPaymentPageRequest struct {
	RazorpayOrderID string `json:"razorpay_order_id" binding:"required"`
}

type ShowPaymentPageResponse struct {
	RazorpayOrderID    string `json:"razorpay_order_id"`
	RazorpayKeyID      string `json:"razorpay_key_id"`
	PlanID             string `json:"plan_id"`
	Amount             int64  `json:"amount"`         // Amount in paise for Razorpay
	DisplayAmount      string `json:"display_amount"` // Formatted amount for display
	PlanDurationInDays int    `json:"plan_duration_in_days"`
	GatewayURL         string `json:"gateway_url"`
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

type PaymentPageData struct {
	OrderID            string  `json:"order_id"`
	RazorpayOrderID    string  `json:"razorpay_order_id"`
	RazorpayKeyID      string  `json:"razorpay_key_id"`
	Amount             float64 `json:"amount"`
	DisplayAmount      string  `json:"display_amount"`
	Currency           string  `json:"currency"`
	PlanID             string  `json:"plan_id"`
	PlanDurationInDays int     `json:"plan_duration_in_days"`
	GatewayURL         string  `json:"gateway_url"`
}

type PaymentSuccessData struct {
	SubscriptionID        string `json:"subscription_id"`
	SubscriptionStartDate string `json:"subscription_start_date"`
	SubscriptionEndDate   string `json:"subscription_end_date"`
	SubscriptionStatus    string `json:"subscription_status"`
}

type PaymentFailureData struct {
	OrderID string `json:"order_id"`
	Error   string `json:"error"`
}

type CreateSubscriptionPlanRequest struct {
	ID           string  `json:"id" binding:"required"`
	DurationDays int     `json:"duration_days" binding:"required"`
	Amount       float64 `json:"amount" binding:"required"`
	Currency     string  `json:"currency" binding:"required"`
	Description  string  `json:"description" binding:"required"`
	IsActive     bool    `json:"is_active" binding:"required"`
}

type UpdateSubscriptionPlanRequest struct {
	ID           string   `json:"id" binding:"required"`
	DurationDays *int     `json:"duration_days,omitempty"`
	Amount       *float64 `json:"amount,omitempty"`
	Currency     *string  `json:"currency,omitempty"`
	Description  *string  `json:"description,omitempty"`
	IsActive     *bool    `json:"is_active,omitempty"`
}

type CreateOrUpdateSubscriptionPlanResponse struct {
	Success bool `json:"success"`
}

type SubscriptionPlan struct {
	ID           string    `json:"id"`
	DurationDays int       `json:"duration_days"`
	Amount       float64   `json:"amount"`
	Currency     string    `json:"currency"`
	Description  string    `json:"description"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type ActiveSubscription struct {
	SubscriptionID int64     `json:"subscription_id"`
	PlanID         string    `json:"plan_id"`
	StartDate      time.Time `json:"start_date"`
	EndDate        time.Time `json:"end_date"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
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
