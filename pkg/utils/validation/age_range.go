package validation

func IsValidAge(age *int) bool {
	if age != nil {
		if *age < 18 || *age > 100 {
			return false
		}
	}
	return true
}

func IsValidAgeRange(minAge, maxAge *int) bool {
	if minAge != nil {
		if *minAge < 18 || *minAge > 100 {
			return false
		}
	}
	if maxAge != nil {
		if *maxAge < 18 || *maxAge > 100 {
			return false
		}
	}
	if minAge != nil && maxAge != nil {
		if *minAge > *maxAge {
			return false
		}
	}
	return true
}
