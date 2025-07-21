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
	Client   *mongo.Client
	Database *mongo.Database
}

func NewClient(ctx context.Context, cfg Config) (*Client, error) {
	if cfg.Timeout == 0 {
		cfg.Timeout = 10 * time.Second
	}

	connCtx, cancel := context.WithTimeout(ctx, cfg.Timeout)
	defer cancel()

	mongoOpts := options.Client().ApplyURI(cfg.URI)
	client, err := mongo.Connect(connCtx, mongoOpts)
	if err != nil {
		return nil, fmt.Errorf("mongodb: connect error: %w", err)
	}

	pingCtx, pingCancel := context.WithTimeout(ctx, 5*time.Second)
	defer pingCancel()
	if err := client.Ping(pingCtx, readpref.Primary()); err != nil {
		_ = client.Disconnect(ctx)
		return nil, fmt.Errorf("mongodb: ping error: %w", err)
	}

	return &Client{
		Client:   client,
		Database: client.Database(cfg.Database),
	}, nil
}

func (c *Client) Close() error {
	if c.Client == nil {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return c.Client.Disconnect(ctx)
}
