package payment

import "time"

type PaymentVerified struct {
	UserID              string    `json:"user_id"`
	SubscriptionID      string    `json:"subscription_id"`
	SubscriptionEndDate time.Time `json:"subscription_end_date"`
	PlanID              string    `json:"plan_id"`
	PaymentID           string    `json:"payment_id"`
	Timestamp           time.Time `json:"timestamp"`
}





