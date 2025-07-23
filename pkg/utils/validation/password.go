package validation

import (
	"unicode"
)

type PasswordRequirements struct {
	MinLength      int
	MaxLength      int
	RequireUpper   bool
	RequireLower   bool
	RequireDigit   bool
	RequireSpecial bool
}

func DefaultPasswordRequirements() PasswordRequirements {
	return PasswordRequirements{
		MinLength:      8,
		MaxLength:      128,
		RequireUpper:   true,
		RequireLower:   true,
		RequireDigit:   true,
		RequireSpecial: true,
	}
}

func IsValidPassword(password string, requirements PasswordRequirements) bool {
	if len(password) < requirements.MinLength || len(password) > requirements.MaxLength {
		return false
	}

	var hasUpper, hasLower, hasDigit, hasSpecial bool

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if requirements.RequireUpper && !hasUpper {
		return false
	}
	if requirements.RequireLower && !hasLower {
		return false
	}
	if requirements.RequireDigit && !hasDigit {
		return false
	}
	if requirements.RequireSpecial && !hasSpecial {
		return false
	}

	return true
}
