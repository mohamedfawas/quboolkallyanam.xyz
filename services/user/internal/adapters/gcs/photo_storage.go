package gcs

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	gcsstore "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/mediastorage/gcs"
	mediastorage "github.com/mohamedfawas/quboolkallyanam.xyz/services/user/internal/domain/mediastorage"
)

type photoStorageAdapter struct {
	gcsStore *gcsstore.GCSStore
	baseURL  string // emulator endpoint like http://localhost:4443 (empty in prod)
	bucket   string
  }

func NewPhotoStorageAdapter(gcsStore *gcsstore.GCSStore, baseURL, bucket string) mediastorage.PhotoStorage {
	return &photoStorageAdapter{
		gcsStore: gcsStore,
		baseURL:  strings.TrimRight(baseURL, "/"),
		bucket:   bucket,
}
}

func (a *photoStorageAdapter) objectURL(key string) string {
	return fmt.Sprintf("%s/%s/%s", a.baseURL, a.bucket, key)
}

func (a *photoStorageAdapter) usingEmulator() bool { return a.baseURL != "" }


func (a *photoStorageAdapter) GetProfilePhotoUploadURL(
	ctx context.Context,
	userID string,
	contentType string,
	expiry time.Duration) (*mediastorage.PhotoUploadURLResponse, error) {
  
	objectKey := fmt.Sprintf("%s/%s", constants.ProfilePhotoStorageDirectory, userID)
  
	var uploadURL string
	var err error
	if a.usingEmulator() {
	  uploadURL = a.objectURL(objectKey)
	} else {
	  uploadURL, err = a.gcsStore.GetUploadURL(objectKey, contentType, expiry)
	  if err != nil {
		return nil, err
	  }
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
  
	var uploadURL string
	var err error
	if a.usingEmulator() {
	  uploadURL = a.objectURL(objectKey)
	} else {
	  uploadURL, err = a.gcsStore.GetUploadURL(objectKey, contentType, expiry)
	  if err != nil {
		return nil, err
	  }
	}
  
	return &mediastorage.PhotoUploadURLResponse{
	  UploadURL:        uploadURL,
	  ObjectKey:        objectKey,
	  ExpiresInSeconds: uint32(expiry.Seconds()),
	}, nil
  }

  func (a *photoStorageAdapter) GetDownloadURL(ctx context.Context, objectKey string, expiry time.Duration) (string, error) {
	if a.usingEmulator() {
	  return a.objectURL(objectKey), nil
	}
	return a.gcsStore.GetDownloadURL(objectKey, expiry)
  }

  func (a *photoStorageAdapter) DeletePhoto(ctx context.Context, objectKey string) error {
	return a.gcsStore.Delete(ctx, objectKey)
  }
