package validation

import "time"

func IsValidDateOfBirth(dateOfBirth string) bool {
	layout := "2006-01-02"
	dob, err := time.Parse(layout, dateOfBirth)
	if err != nil {
		return false
	}

	now := time.Now()
	minDate := now.AddDate(-100, 0, 0) // Not older than 100 years
	maxDate := now.AddDate(-18, 0, 0)  // Must be at least 18 years old

	return dob.After(minDate) && dob.Before(maxDate)
}
