package validation

// If you are adding more types, please update it in ErrInvalidImageType public message as well
var allowedImageType = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/jpg": true,
}

func IsValidImageType(contentType string) bool {
	return allowedImageType[contentType]
}