package firestore

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
)

type Config struct {
	ProjectID       string
	CredentialsFile string
	EmulatorHost    string
}

type Client struct {
	client *firestore.Client
}

// NewClient creates a new Firestore client following industry standard pattern
func NewClient(ctx context.Context, projectID string, opts ...option.ClientOption) (*Client, error) {
	client, err := firestore.NewClient(ctx, projectID, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create firestore client: %w", err)
	}

	return &Client{
		client: client,
	}, nil
}

// NewClientWithConfig creates a new Firestore client with custom config
func NewClientWithConfig(ctx context.Context, cfg Config) (*Client, error) {
	var opts []option.ClientOption

	if cfg.CredentialsFile != "" {
		opts = append(opts, option.WithCredentialsFile(cfg.CredentialsFile))
	}

	return NewClient(ctx, cfg.ProjectID, opts...)
}

func (c *Client) Close() error {
	if c.client != nil {
		return c.client.Close()
	}
	return nil
}

// GetClient returns the underlying Firestore client for advanced operations
func (c *Client) GetClient() *firestore.Client {
	return c.client
}
