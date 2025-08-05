package mediastorage

import (
	"context"
	"time"
)

type PhotoUploadURLResponse struct {
	UploadURL        string
	ObjectKey        string
	ExpiresInSeconds uint32
}

type PhotoStorage interface {
	GetProfilePhotoUploadURL(
		ctx context.Context,
		userID string,
		contentType string,
		expiry time.Duration) (*PhotoUploadURLResponse, error)
	GetAdditionalPhotoUploadURL(
		ctx context.Context,
		userID string,
		displayOrder int32,
		contentType string,
		expiry time.Duration) (*PhotoUploadURLResponse, error)
	GetDownloadURL(ctx context.Context, objectKey string, expiry time.Duration) (string, error)
	DeletePhoto(ctx context.Context, objectKey string) error
}
