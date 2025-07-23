package validation

import (
	"regexp"
	"strings"
)

const (
	minFullNameLength = 2
	maxFullNameLength = 100
)

var fullNameRegex = regexp.MustCompile(`^[A-Za-z]+(?: [A-Za-z]+)*$`)

func IsValidFullName(name string) bool {
	name = strings.TrimSpace(name)
	length := len(name)

	if length < minFullNameLength || length > maxFullNameLength {
		return false
	}

	return fullNameRegex.MatchString(name)
}
