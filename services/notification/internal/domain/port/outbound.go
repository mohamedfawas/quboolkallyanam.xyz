package port

import (
	"context"

	"github.com/mohamedfawas/quboolkallyanam.xyz/services/notification/internal/domain/model"
)

type EmailAdapter interface {
	SendEmail(ctx context.Context, req model.EmailRequest) error
}
