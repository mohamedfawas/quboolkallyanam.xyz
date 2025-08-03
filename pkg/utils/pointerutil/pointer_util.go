package pointerutil

func GetStringValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func GetIntValue(val *int) int {
    if val != nil {
        return *val
    }
    return 0 
}

func GetBoolValue(value *bool) bool {
	if value == nil {
		return false
	}
	return *value
}