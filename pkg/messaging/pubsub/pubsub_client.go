package pubsub

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/messaging"

	"cloud.google.com/go/pubsub"
)

type Client struct {
	psClient *pubsub.Client
	project  string
	ctx      context.Context
}

var _ messaging.Client = (*Client)(nil)

func NewClient(ctx context.Context, projectID string) (*Client, error) {
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to create pubsub client: %w", err)
	}
	return &Client{
		psClient: client,
		project:  projectID,
		ctx:      ctx,
	}, nil
}

func (c *Client) Publish(topicName string, data interface{}) error {
	topic := c.psClient.Topic(topicName)

	bytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	ctx, cancel := context.WithTimeout(c.ctx, 5*time.Second)
	defer cancel()

	result := topic.Publish(ctx, &pubsub.Message{
		Data: bytes,
	})

	_, err = result.Get(ctx)
	return err
}

func (c *Client) Subscribe(subName string, handler messaging.MessageHandler) error {
	sub := c.psClient.Subscription(subName)

	go func() {
		err := sub.Receive(c.ctx, func(ctx context.Context, msg *pubsub.Message) {
			if err := handler(msg.Data); err != nil {
				log.Printf("❌ Handler error: %v", err)
				msg.Nack()
			} else {
				msg.Ack()
			}
		})
		if err != nil {
			log.Printf("❌ Subscription error: %v", err)
		}
	}()

	return nil
}

func (c *Client) Close() error {
	return c.psClient.Close()
}
