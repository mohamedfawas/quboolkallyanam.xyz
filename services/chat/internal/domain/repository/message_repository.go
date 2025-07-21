package repository

import (
	"context"

	"github.com/mohamedfawas/quboolkallyanam.xyz/services/chat/internal/domain/entity"
)

type MessageRepository interface {
	CreateMessage(ctx context.Context, message *entity.Message) error
}
