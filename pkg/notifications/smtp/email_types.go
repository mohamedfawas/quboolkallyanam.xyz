package smtp

type EmailType string

const (
	EmailTypeOTPVerification  EmailType = "otp_verification"
	EmailTypeInterestReceived EmailType = "interest_received"
	EmailTypePaymentSuccess     EmailType = "payment_success"
)
