package smtp

import (
	"bytes"
	"fmt"
	"net/smtp"

	"errors"
)

type Config struct {
	SMTPHost     string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string
	FromEmail    string
	FromName     string
}

type Client struct {
	config Config
	auth   smtp.Auth
	addr   string
}

func NewClient(cfg Config) (*Client, error) {
	if cfg.SMTPHost == "" || cfg.SMTPPort == 0 || cfg.FromEmail == "" {
		return nil, errors.New("smtp: invalid config, missing host/port/from")
	}
	// For real servers, require credentials; MailHog doesn’t enforce auth
	var auth smtp.Auth
	if cfg.SMTPHost != "mailhog" {
		if cfg.SMTPUsername == "" || cfg.SMTPPassword == "" {
			return nil, errors.New("smtp: invalid config, missing credentials")
		}
		auth = smtp.PlainAuth("", cfg.SMTPUsername, cfg.SMTPPassword, cfg.SMTPHost)
	}

	return &Client{
		config: cfg,
		auth:   auth,
		addr:   fmt.Sprintf("%s:%d", cfg.SMTPHost, cfg.SMTPPort),
	}, nil
}

func (c *Client) SendEmail(to, subject, body string) error {
	// Build the “From” header with optional FromName
	from := c.config.FromEmail
	if c.config.FromName != "" {
		from = fmt.Sprintf("%s <%s>", c.config.FromName, c.config.FromEmail)
	}

	// Build headers
	var msg bytes.Buffer
	msg.WriteString(fmt.Sprintf("From: %s\r\n", from))
	msg.WriteString(fmt.Sprintf("To: %s\r\n", to))
	msg.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	// If you ever want HTML support, switch Content‑Type to text/html
	msg.WriteString("MIME-Version: 1.0\r\n")
	msg.WriteString("Content-Type: text/plain; charset=\"utf-8\"\r\n")
	msg.WriteString("\r\n") // blank line between headers and body
	msg.WriteString(body)   // the email body

	// Actually send
	return smtp.SendMail(
		c.addr,             // "smtp.example.com:587"
		c.auth,             // AUTH if configured (nil for MailHog)
		c.config.FromEmail, // envelope from
		[]string{to},       // envelope to
		msg.Bytes(),        // full message
	)
}
