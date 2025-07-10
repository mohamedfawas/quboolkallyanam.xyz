package smtp

import (
	"fmt"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/notifications/smtp/templates"
)

type EmailRequest struct {
	To      string
	Type    EmailType
	Subject string
	Payload map[string]string
}

func (c *Client) SendEmailByType(req EmailRequest) error {
	var body string

	switch req.Type {
	case EmailTypeOTPVerification:
		body = templates.BuildOTPBody(req.Payload["email"], req.Payload["otp"], req.Payload["expiryMinutes"])
	// case EmailTypeInterestReceived:
	// 	body = templates.BuildInterestReceivedBody(req.Payload["email"], req.Payload["interest"])
	default:
		return fmt.Errorf("unknown email type: %s", req.Type)
	}

	return c.SendEmail(req.To, req.Subject, body)
}
