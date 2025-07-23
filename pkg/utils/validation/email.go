package validation

import (
	"regexp"
	"strings"
)

func IsValidEmail(email string) bool {
	email = strings.TrimSpace(email)

	if email == "" {
		return false
	}

	// Basic length check (email addresses shouldn't be extremely long)
	if len(email) > 254 {
		return false
	}

	// Regular expression pattern for email validation
	// This pattern covers most common email formats
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	// Compile the regex
	re := regexp.MustCompile(emailRegex)

	// Test the email against the pattern
	return re.MatchString(email)
}
