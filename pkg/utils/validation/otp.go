package validation

import (
	"regexp"
	"strings"
)

func IsValidOTP(otp string, length int) bool {
	otp = strings.TrimSpace(otp)

	if otp == "" {
		return false
	}

	if len(otp) != length {
		return false
	}

	digitRegex := `^[0-9]+$`
	re := regexp.MustCompile(digitRegex)

	return re.MatchString(otp)
}
