package validation

func IsValidHumanHeight(height int) bool {
	return height >= 100 && height <= 250
}
