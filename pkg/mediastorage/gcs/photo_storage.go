package gcs

import (
	"context"
	"fmt"
	"os"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

type MediaStorageConfig struct {
	Bucket          string        // GCS bucket name (required)
	CredentialsFile string        // path to service‐account JSON (only for real GCS)
	SignerEmail     string        // service account email for SignedURL (optional)
	PrivateKeyPath  string        // path to PEM for SignedURL (optional)
	URLExpiry       time.Duration // default expiry for signed URLs
	Endpoint        string        // emulator endpoint, e.g. "http://localhost:4443"; empty → real GCS
}

type GCSStore struct {
	client      *storage.Client
	bucket      string
	signerEmail string // service account email for SignedURL
	privateKey  []byte // PEM bytes
}

func NewGCSStore(
	ctx context.Context,
	config MediaStorageConfig,
) (*GCSStore, error) {
	if config.Bucket == "" {
		return nil, fmt.Errorf("bucket name is required")
	}

	// Build client options
	opts := []option.ClientOption{}
	if config.Endpoint != "" {
		// DEV: hit emulator, no auth
		opts = append(opts,
			option.WithEndpoint(config.Endpoint),
			option.WithoutAuthentication(),
		)
	} else {
		// PROD: real GCS, load creds file if provided
		if config.CredentialsFile != "" {
			opts = append(opts,
				option.WithCredentialsFile(config.CredentialsFile),
			)
		}
	}

	cli, err := storage.NewClient(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("gcs client: %w", err)
	}

	// Load private key if path provided (for Signed URLs)
	var pk []byte
	if config.PrivateKeyPath != "" {
		pk, err = os.ReadFile(config.PrivateKeyPath)
		if err != nil {
			cli.Close()
			return nil, fmt.Errorf("read private key: %w", err)
		}
	}

	// If you’ve loaded a key, require a signer email
	if len(pk) > 0 && config.SignerEmail == "" {
		cli.Close()
		return nil, fmt.Errorf("signer email is required when private key is provided")
	}	

	return &GCSStore{
		client:      cli,
		bucket:      config.Bucket,
		signerEmail: config.SignerEmail,
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
