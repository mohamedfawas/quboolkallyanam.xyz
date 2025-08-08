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

const MaxImageFileSize = 5 * 1024 * 1024 // 5MB

func IsValidImageFileSize(fileSize int64) bool {
	return fileSize <= MaxImageFileSize && fileSize > 0
}
