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
}

func NewPhotoStorageAdapter(gcsStore *gcsstore.GCSStore) mediastorage.PhotoStorage {
	return &photoStorageAdapter{
		gcsStore: gcsStore,
	}
}

func (a *photoStorageAdapter) GetProfilePhotoUploadURL(
	ctx context.Context,
	userID string,
	contentType string,
	expiry time.Duration) (*mediastorage.PhotoUploadURLResponse, error) {
	objectKey := fmt.Sprintf("%s/%s", constants.ProfilePhotoStorageDirectory, userID)

	uploadURL, err := a.gcsStore.GetUploadURL(objectKey, contentType, expiry)
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

	uploadURL, err := a.gcsStore.GetUploadURL(objectKey, contentType, expiry)
	if err != nil {
		return nil, err
	}

	return &mediastorage.PhotoUploadURLResponse{
		UploadURL:        uploadURL,
		ObjectKey:        objectKey,
		ExpiresInSeconds: uint32(expiry.Seconds()),
	}, nil
}

func (p *photoStorageAdapter) GetDownloadURL(ctx context.Context, objectKey string, expiry time.Duration) (string, error) {
	downloadURL, err := p.gcsStore.GetDownloadURL(objectKey, expiry)
	if err != nil {
		return "", fmt.Errorf("failed to generate download URL: %w", err)
	}
	return downloadURL, nil
}

func (p *photoStorageAdapter) DeletePhoto(ctx context.Context, objectKey string) error {
	if err := p.gcsStore.Delete(ctx, objectKey); err != nil {
		return fmt.Errorf("failed to delete photo: %w", err)
	}
	return nil
}
