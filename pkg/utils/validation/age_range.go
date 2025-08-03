package validation

func IsValidAge(age int) bool {
	if age < 18 || age > 100 {
		return false
	}
	return true
}

func IsValidAgeRange(minAge, maxAge int) bool {
	if minAge < 18 || minAge > 100 {
		return false
	}
	if maxAge < 18 || maxAge > 100 {
		return false
	}
	if minAge > maxAge {
		return false
	}
	return true
}
