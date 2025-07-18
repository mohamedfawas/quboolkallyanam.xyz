package model

type EmailType string

const (
	EmailTypeOTPVerification  EmailType = "otp_verification"
	EmailTypePaymentSuccess   EmailType = "payment_success"
	EmailTypeInterestReceived EmailType = "interest_received"
)