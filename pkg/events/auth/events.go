package authevents

type UserOTPRequestedEvent struct {
	Email         string `json:"email"`
	OTP           string `json:"otp"`
	ExpiryMinutes int    `json:"expiry_minutes"`
}
