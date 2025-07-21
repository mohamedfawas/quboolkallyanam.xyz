package usecase

import (
	"context"

	"github.com/google/uuid"
)

type UserProfileUsecase interface {
	UpdateUserLastLogin(ctx context.Context,
		userID uuid.UUID,
		email, phone string) error
}
