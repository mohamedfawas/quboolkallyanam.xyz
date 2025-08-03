package validation

func IsValidHeight(height int) bool {
	if height < 130 || height > 220 {
		return false
	}
	return true
}

func IsValidHeightRange(minHeight, maxHeight int) bool {
	if minHeight < 130 || minHeight > 220 {
		return false
	}
	if maxHeight < 130 || maxHeight > 220 {
		return false
	}
	if minHeight > maxHeight {
		return false
	}
	return true
}
