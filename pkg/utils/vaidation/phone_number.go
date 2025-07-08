package validation

import (
	"regexp"
	"strings"
)

func IsValidPhoneNumber(phone string) bool {
	// Remove all whitespace and common separators
	phone = strings.ReplaceAll(phone, " ", "")
	phone = strings.ReplaceAll(phone, "-", "")
	phone = strings.ReplaceAll(phone, "(", "")
	phone = strings.ReplaceAll(phone, ")", "")
	phone = strings.ReplaceAll(phone, ".", "")

	if phone == "" {
		return false
	}

	// Remove country code prefix if present
	if strings.HasPrefix(phone, "+") {
		phone = phone[1:]
	}

	// Check length (should be between 7-15 digits for most international numbers)
	if len(phone) < 7 || len(phone) > 15 {
		return false
	}

	// Check if all remaining characters are digits
	digitRegex := `^[0-9]+$`
	re := regexp.MustCompile(digitRegex)

	return re.MatchString(phone)
}
