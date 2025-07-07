package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/messaging"
	amqp "github.com/rabbitmq/amqp091-go"
)

var _ messaging.Client = (*Client)(nil)

type Config struct {
	DSN          string // amqp://guest:guest@localhost:5672/
	ExchangeName string // e.g., "events"
}

type Client struct {
	conn     *amqp.Connection
	channel  *amqp.Channel
	exchange string
}

func NewClient(cfg Config) (*Client, error) {
	conn, err := amqp.Dial(cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	err = ch.ExchangeDeclare(
		cfg.ExchangeName,
		"topic",
		true,  // durable
		false, // auto-deleted
		false,
		false,
		nil,
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare exchange: %w", err)
	}

	log.Println("✅ Connected to RabbitMQ")
	return &Client{
		conn:     conn,
		channel:  ch,
		exchange: cfg.ExchangeName,
	}, nil
}

func (c *Client) Publish(routingKey string, message interface{}) error {
	payload, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return c.channel.PublishWithContext(
		ctx,
		c.exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        payload,
		},
	)
}

func (c *Client) Subscribe(routingKey string, handler messaging.MessageHandler) error {
	q, err := c.channel.QueueDeclare(
		"",    // auto-generated queue name
		false, // non-durable
		true,  // auto-delete
		true,  // exclusive
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("queue declare failed: %w", err)
	}

	err = c.channel.QueueBind(
		q.Name,
		routingKey,
		c.exchange,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("queue bind failed: %w", err)
	}

	msgs, err := c.channel.Consume(
		q.Name,
		"",
		true,  // auto-ack
		false, // not exclusive
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("consume failed: %w", err)
	}

	go func() {
		for d := range msgs {
			if err := handler(d.Body); err != nil {
				log.Printf("❌ Error in handler: %v", err)
			}
		}
	}()

	return nil
}

func (c *Client) Close() error {
	if c.channel != nil {
		_ = c.channel.Close()
	}
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
