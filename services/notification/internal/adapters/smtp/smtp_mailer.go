package smtp

import (
	"context"
	"fmt"
	"log"

	pkgSMTP "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/notifications/smtp"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/notification/internal/domain/model"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/notification/internal/domain/port"
)

type emailAdapter struct {
	client *pkgSMTP.Client
}

func NewEmailAdapter(client *pkgSMTP.Client) port.EmailAdapter {
	return &emailAdapter{
		client: client,
	}
}

func (e *emailAdapter) SendEmail(ctx context.Context, req model.EmailRequest) error {
	if err := e.client.SendEmail(req.To, req.Subject, req.Body); err != nil {
		log.Printf("smtp adapter: failed to send email: %v", err)
		return fmt.Errorf("smtp adapter: failed to send email: %w", err)
	}
	return nil
}
