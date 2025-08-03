package gcs

import (
	"context"
	"fmt"
	"os"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

type GCSStore struct {
	client      *storage.Client
	bucket      string
	signerEmail string // service account email for SignedURL
	privateKey  []byte // PEM bytes
}

func NewGCSStore(
	ctx context.Context,
	bucket string,
	credsFile string,
	signerEmail string,
	privateKeyPath string,
) (*GCSStore, error) {
	if bucket == "" {
		return nil, fmt.Errorf("bucket name is required")
	}

	opts := []option.ClientOption{}
	if credsFile != "" {
		opts = append(opts, option.WithCredentialsFile(credsFile))
	}

	cli, err := storage.NewClient(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("gcs client: %w", err)
	}

	var pk []byte
	if privateKeyPath != "" {
		pk, err = os.ReadFile(privateKeyPath)
		if err != nil {
			return nil, fmt.Errorf("read private key: %w", err)
		}
	}

	// Validate signer email if private key is provided
	if len(pk) > 0 && signerEmail == "" {
		return nil, fmt.Errorf("signer email is required when private key is provided")
	}

	return &GCSStore{
		client:      cli,
		bucket:      bucket,
		signerEmail: signerEmail,
		privateKey:  pk,
	}, nil
}

func (s *GCSStore) GetUploadURL(
	key, contentType string,
	expiry time.Duration,
) (string, error) {
	if len(s.privateKey) == 0 || s.signerEmail == "" {
		return "", fmt.Errorf("private key and signer email are required for signed URLs")
	}

	opts := &storage.SignedURLOptions{
		GoogleAccessID: s.signerEmail,
		PrivateKey:     s.privateKey,
		Method:         "PUT",
		Expires:        time.Now().Add(expiry),
		ContentType:    contentType,
	}

	url, err := storage.SignedURL(s.bucket, key, opts)
	if err != nil {
		return "", fmt.Errorf("signed PUT URL: %w", err)
	}
	return url, nil
}

func (s *GCSStore) GetDownloadURL(
	key string,
	expiry time.Duration,
) (string, error) {
	if len(s.privateKey) == 0 || s.signerEmail == "" {
		return "", fmt.Errorf("private key and signer email are required for signed URLs")
	}

	opts := &storage.SignedURLOptions{
		GoogleAccessID: s.signerEmail,
		PrivateKey:     s.privateKey,
		Method:         "GET",
		Expires:        time.Now().Add(expiry),
	}

	url, err := storage.SignedURL(s.bucket, key, opts)
	if err != nil {
		return "", fmt.Errorf("signed GET URL: %w", err)
	}
	return url, nil
}

func (s *GCSStore) Delete(ctx context.Context, key string) error {
	return s.client.Bucket(s.bucket).Object(key).Delete(ctx)
}

func (s *GCSStore) Close() error {
	return s.client.Close()
}
