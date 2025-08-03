package ageutil


import (
	"time"
)


func CalculateAge(dob time.Time) int {
	now := time.Now().UTC()
	age := now.Year() - dob.Year()

	if now.Month() < dob.Month() || (now.Month() == dob.Month() && now.Day() < dob.Day()) {
		age--
	}

	return age
}