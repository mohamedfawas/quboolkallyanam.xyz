package gcs

import (
	"context"
	"fmt"
	"time"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	gcsstore "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/mediastorage/gcs"
	mediastorage "github.com/mohamedfawas/quboolkallyanam.xyz/services/user/internal/domain/mediastorage"
)

type photoStorageAdapter struct {
	gcsStore *gcsstore.GCSStore
	bucket   string
}

func NewPhotoStorageAdapter(gcsStore *gcsstore.GCSStore, bucket string) mediastorage.PhotoStorage {
	return &photoStorageAdapter{
		gcsStore: gcsStore,
		bucket:   bucket,
	}
}

func (a *photoStorageAdapter) GetProfilePhotoUploadURL(
	ctx context.Context,
	userID string,
	contentType string,
	expiry time.Duration) (*mediastorage.PhotoUploadURLResponse, error) {

	objectKey := fmt.Sprintf("%s/%s", constants.ProfilePhotoStorageDirectory, userID)

	uploadURL, err := a.gcsStore.GetUploadURL(ctx, objectKey, contentType, expiry)
	if err != nil {
		return nil, err
	}

	return &mediastorage.PhotoUploadURLResponse{
		UploadURL:        uploadURL,
		ObjectKey:        objectKey,
		ExpiresInSeconds: uint32(expiry.Seconds()),
	}, nil
}



func (a *photoStorageAdapter) GetAdditionalPhotoUploadURL(
	ctx context.Context,
	userID string,
	displayOrder int32,
	contentType string,
	expiry time.Duration) (*mediastorage.PhotoUploadURLResponse, error) {

	objectKey := fmt.Sprintf("%s/%s/%d", constants.AdditionalPhotoStorageDirectory, userID, displayOrder)

	uploadURL, err := a.gcsStore.GetUploadURL(ctx, objectKey, contentType, expiry)
	if err != nil {
		return nil, err
	}

	return &mediastorage.PhotoUploadURLResponse{
		UploadURL:        uploadURL,
		ObjectKey:        objectKey,
		ExpiresInSeconds: uint32(expiry.Seconds()),
	}, nil
}

func (a *photoStorageAdapter) GetDownloadURL(ctx context.Context, objectKey string, expiry time.Duration) (string, error) {
	return a.gcsStore.GetDownloadURL(ctx, objectKey, expiry)
}

func (a *photoStorageAdapter) DeletePhoto(ctx context.Context, objectKey string) error {
	return a.gcsStore.Delete(ctx, objectKey)
}
