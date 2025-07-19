package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type UserProfileUsecase interface {
	UpdateUserLastLogin(ctx context.Context,
		userID uuid.UUID,
		email, phone string) error
}
