package validation

func IsValidHeight(height *int) bool {
	if height == nil {
		return false
	}
	return *height >= 130 && *height <= 220
}

func IsValidHeightRange(minHeight, maxHeight *int) bool {
	if minHeight != nil {
		if *minHeight < 130 || *minHeight > 220 {
			return false
		}
	}
	if maxHeight != nil {
		if *maxHeight < 130 || *maxHeight > 220 {
			return false
		}
	}
	if minHeight != nil && maxHeight != nil {
		if *minHeight > *maxHeight {
			return false
		}
	}
	return true
}
