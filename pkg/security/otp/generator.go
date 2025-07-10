package otp

import (
	"crypto/rand"
	"fmt"
)

func GenerateNumericOTP(length int) (string, error) {
	const digits = "0123456789"
	otp := make([]byte, length)

	_, err := rand.Read(otp)
	if err != nil {
		return "", fmt.Errorf("failed to generate otp: %w", err)
	}

	for i := range otp {
		otp[i] = digits[otp[i]%10]
	}

	return string(otp), nil
}
