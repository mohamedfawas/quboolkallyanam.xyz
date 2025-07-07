package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Config struct {
	URI      string
	Database string
	Timeout  time.Duration
}

type Client struct {
	client   *mongo.Client
	database *mongo.Database
}

// NewClient creates a new MongoDB client following industry standard pattern
func NewClient(ctx context.Context, cfg Config) (*Client, error) {
	// Set default timeout if not provided
	if cfg.Timeout == 0 {
		cfg.Timeout = 10 * time.Second
	}

	// Create connection context with timeout
	connCtx, cancel := context.WithTimeout(ctx, cfg.Timeout)
	defer cancel()

	// Create client options
	opts := options.Client().ApplyURI(cfg.URI)

	// Connect to MongoDB
	client, err := mongo.Connect(connCtx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mongodb: %w", err)
	}

	// Ping to verify connection
	pingCtx, pingCancel := context.WithTimeout(ctx, 5*time.Second)
	defer pingCancel()

	if err := client.Ping(pingCtx, readpref.Primary()); err != nil {
		client.Disconnect(ctx) // Clean up on ping failure
		return nil, fmt.Errorf("failed to ping mongodb: %w", err)
	}

	database := client.Database(cfg.Database)

	return &Client{
		client:   client,
		database: database,
	}, nil
}

func (c *Client) Close() error {
	if c.client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return c.client.Disconnect(ctx)
	}
	return nil
}

// GetClient returns the underlying MongoDB client for advanced operations
func (c *Client) GetClient() *mongo.Client {
	return c.client
}

// GetDatabase returns the database instance
func (c *Client) GetDatabase() *mongo.Database {
	return c.database
}
