package gcs

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"cloud.google.com/go/compute/metadata"
	"cloud.google.com/go/iam/credentials/apiv1"
	"cloud.google.com/go/storage"
	credentialspb "google.golang.org/genproto/googleapis/iam/credentials/v1"
	"google.golang.org/api/googleapi"
)
/*
this code was used in development,  because we were using fake-gcs-server for testing.
======================================================================================
type MediaStorageConfig struct {
	Bucket          string        // GCS bucket name (required)
	CredentialsFile string        // path to service‐account JSON (only for real GCS)
	SignerEmail     string        // service account email for SignedURL (optional)
	PrivateKeyPath  string        // path to PEM for SignedURL (optional)
	URLExpiry       time.Duration // default expiry for signed URLs
	Endpoint        string        // emulator endpoint, e.g. "http://localhost:4443"; empty → real GCS
}
======================================================================================
*/

type MediaStorageConfig struct {
	Bucket      string        
	URLExpiry   time.Duration // default expiry for signed URLs
	SignerEmail string        // service account email to sign URLs (recommended)
}

type GCSStore struct {
	client        *storage.Client
	bucket        string
	defaultExpiry time.Duration

	// signer email configured via MediaStorageConfig.SignerEmail (or auto-discovered on GCP)
	signerEmail string

	// iamClient talks to the IAM Credentials API; created lazily when signing is needed
	iamClient *credentials.IamCredentialsClient
	mu        sync.Mutex // protects iamClient
}


// NewGCSStore initializes the storage client and records signer email (if provided).
// If SignerEmail is empty and running on GCP, it will try to auto-discover the default SA email.
func NewGCSStore(ctx context.Context, cfg MediaStorageConfig) (*GCSStore, error) {
	if cfg.Bucket == "" {
		return nil, fmt.Errorf("bucket name is required")
	}
	if cfg.URLExpiry == 0 {
		cfg.URLExpiry = 15 * time.Minute
	}

	cli, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("storage.NewClient: %w. set GOOGLE_APPLICATION_CREDENTIALS or run on GCP; see https://cloud.google.com/docs/authentication/external/set-up-adc", err)
	}

	signer := cfg.SignerEmail
	if signer == "" && metadata.OnGCE() {
		if email, err := metadata.GetWithContext(ctx, "instance/service-accounts/default/email"); err == nil {
			signer = email
		}
	}

	return &GCSStore{
		client:        cli,
		bucket:        cfg.Bucket,
		defaultExpiry: cfg.URLExpiry,
		signerEmail:   signer,
	}, nil
}


// obtainIAMClient creates the IAM client on first use, concurrency-safe.
func (s *GCSStore) obtainIAMClient(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.iamClient != nil {
		return nil
	}
	if s.signerEmail == "" {
		return fmt.Errorf("signer email not configured")
	}
	iamCli, err := credentials.NewIamCredentialsClient(ctx)
	if err != nil {
		return fmt.Errorf("creating iam credentials client: %w", err)
	}
	s.iamClient = iamCli
	return nil
}


// GetUploadURL returns a V4 signed PUT URL for direct upload.
// ctx allows caller to control timeouts/cancellation.
func (s *GCSStore) GetUploadURL(ctx context.Context, key, contentType string, expiry time.Duration) (string, error) {
	if expiry == 0 {
		expiry = s.defaultExpiry
	}
	return s.signedURL(ctx, key, "PUT", expiry, contentType)
}

// GetDownloadURL returns a V4 signed GET URL for downloading the object.
func (s *GCSStore) GetDownloadURL(ctx context.Context, key string, expiry time.Duration) (string, error) {
	if expiry == 0 {
		expiry = s.defaultExpiry
	}
	return s.signedURL(ctx, key, "GET", expiry, "")
}


// signedURL is the shared implementation that uses storage.SignedURL with SignBytes via IAM SignBlob.
func (s *GCSStore) signedURL(ctx context.Context, key, method string, expiry time.Duration, contentType string) (string, error) {
	if expiry < 0 {
		return "", fmt.Errorf("expiry must be >= 0")
	}

	// GCS V4 signed URL maximum expiry is 7 days.
	const maxExpiry = 7 * 24 * time.Hour
	if expiry > maxExpiry {
		expiry = maxExpiry
	}

	if s.signerEmail == "" {
		return "", fmt.Errorf("signer email not configured: set MediaStorageConfig.SignerEmail or run on GCP where it can be discovered")
	}

	// ensure iam client exists
	if err := s.obtainIAMClient(ctx); err != nil {
		return "", err
	}

	opts := &storage.SignedURLOptions{
		Scheme:      storage.SigningSchemeV4,
		Method:      method,
		Expires:     time.Now().Add(expiry),
		ContentType: contentType,
		GoogleAccessID: s.signerEmail,
		SignBytes: func(b []byte) ([]byte, error) {
			req := &credentialspb.SignBlobRequest{
				Name:    "projects/-/serviceAccounts/" + s.signerEmail,
				Payload: b,
			}
			// use caller's ctx for IAM call
			resp, err := s.iamClient.SignBlob(ctx, req)
			if err != nil {
				return nil, fmt.Errorf("iam SignBlob: %w", err)
			}
			return resp.SignedBlob, nil
		},
	}

	url, err := storage.SignedURL(s.bucket, key, opts)
	if err != nil {
		return "", fmt.Errorf("storage.SignedURL: %w", err)
	}
	return url, nil
}

func (s *GCSStore) Delete(ctx context.Context, key string) error {
	err := s.client.Bucket(s.bucket).Object(key).Delete(ctx)
	if err == nil {
		return nil
	}
	if errors.Is(err, storage.ErrObjectNotExist) {
		return nil
	}
	var gErr *googleapi.Error
	if errors.As(err, &gErr) && gErr.Code == http.StatusNotFound {
		return nil
	}
	return err
}

func (s *GCSStore) Close() error {
	var firstErr error
	s.mu.Lock()
	if s.iamClient != nil {
		if err := s.iamClient.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
		s.iamClient = nil
	}
	s.mu.Unlock()

	if err := s.client.Close(); err != nil && firstErr == nil {
		firstErr = err
	}
	return firstErr
}
