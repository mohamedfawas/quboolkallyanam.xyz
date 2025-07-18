package templates

import "fmt"

func BuildOTPVerificationBody(email, otp, expiryMinutes string) string {
	return fmt.Sprintf(
		`Hello %s,

Your One-Time Password (OTP) for verifying your account on Qubool Kallyanam is:

	%s

Please enter this OTP to complete your verification. It will expire in %s minutes.

If you didn't request this, you can safely ignore this email.

Regards,  
Team Qubool Kallyanam`,
		email,
		otp,
		expiryMinutes,
	)
}