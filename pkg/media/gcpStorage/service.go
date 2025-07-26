package gcpstorage

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/google/uuid"
	"google.golang.org/api/option"
)

// Error definitions
var (
	ErrInvalidFile     = fmt.Errorf("invalid file")
	ErrFileTooLarge    = fmt.Errorf("file exceeds max size")
	ErrInvalidFileType = fmt.Errorf("file type not supported")
)

// MaxFileSize limits uploads to 5MB
const MaxFileSize = 5 * 1024 * 1024

// AllowedFileTypes lists permitted image extensions
var AllowedFileTypes = []string{".jpg", ".jpeg", ".png"}

// isValidExt checks file extension against AllowedFileTypes
func isValidExt(ext string) bool {
	for _, allow := range AllowedFileTypes {
		if ext == allow {
			return true
		}
	}
	return false
}

// Service wraps a GCS client for photo operations
type Service struct {
	client    *storage.Client
	bucket    string
	urlExpiry time.Duration
}

// Config contains settings for GCS storage
type Config struct {
	Bucket          string        // GCS bucket name
	CredentialsFile string        // path to JSON credentials (leave blank for ADC)
	URLExpiry       time.Duration // signed URL duration
}

// NewService creates a new storage Service
func NewService(ctx context.Context, cfg Config) (*Service, error) {
	opts := []option.ClientOption{}
	if cfg.CredentialsFile != "" {
		opts = append(opts, option.WithCredentialsFile(cfg.CredentialsFile))
	}
	client, err := storage.NewClient(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("gcs NewClient: %w", err)
	}
	return &Service{
		client:    client,
		bucket:    cfg.Bucket,
		urlExpiry: cfg.URLExpiry,
	}, nil
}

// UploadPhoto uploads an image file under folder/userID with a UUID filename
func (s *Service) UploadPhoto(
	ctx context.Context,
	folder, userID string,
	file *multipart.FileHeader,
) (string, error) {
	// Validate size
	if file.Size > MaxFileSize {
		return "", ErrFileTooLarge
	}
	// Validate extension
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !isValidExt(ext) {
		return "", ErrInvalidFileType
	}

	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// Generate storage key: folder/userID/{uuid}.ext
	key := fmt.Sprintf("%s/%s/%s%s", folder, userID, uuid.NewString(), ext)

	w := s.client.Bucket(s.bucket).Object(key).NewWriter(ctx)
	// Determine content type
	w.ContentType = file.Header.Get("Content-Type")
	if w.ContentType == "" {
		buf := make([]byte, 512)
		n, _ := src.Read(buf)
		w.ContentType = http.DetectContentType(buf[:n])
		src.Seek(0, io.SeekStart)
	}

	// Upload data
	if _, err := io.Copy(w, src); err != nil {
		w.Close()
		return "", err
	}
	if err := w.Close(); err != nil {
		return "", err
	}

	// Return public URL; adjust if your bucket is private
	url := fmt.Sprintf("https://storage.googleapis.com/%s/%s", s.bucket, key)
	return url, nil
}

// UploadProfilePhoto is a convenience wrapper for the "profile-photos" folder
func (s *Service) UploadProfilePhoto(ctx context.Context, userID string, file *multipart.FileHeader) (string, error) {
	return s.UploadPhoto(ctx, "profile-photos", userID, file)
}

// UploadAdditionalPhoto is a convenience wrapper for the "additional-photos" folder
func (s *Service) UploadAdditionalPhoto(ctx context.Context, userID string, file *multipart.FileHeader) (string, error) {
	return s.UploadPhoto(ctx, "additional-photos", userID, file)
}

// DeletePhoto removes an object by its key
func (s *Service) DeletePhoto(ctx context.Context, key string) error {
	return s.client.Bucket(s.bucket).Object(key).Delete(ctx)
}

// SignedURL generates a time-limited GET URL for a private bucket
// Requires service account credentials with signing permissions
func (s *Service) SignedURL(key string) (string, error) {
	opts := &storage.SignedURLOptions{
		GoogleAccessID: "<SERVICE-ACCOUNT-EMAIL>",
		PrivateKey:     []byte("<PEM-PRIVATE-KEY>"),
		Method:         "GET",
		Expires:        time.Now().Add(s.urlExpiry),
	}
	return storage.SignedURL(s.bucket, key, opts)
}
