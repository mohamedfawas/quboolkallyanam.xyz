package user

import (
	"context"
	"fmt"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apperrors"
	authevents "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/events/auth"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/security/hash"
)

func (u *userUseCase) UserAccountDelete(ctx context.Context, userID string, password string) error {
	user, err := u.userRepository.GetUser(ctx, "id", userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	if user == nil {
		return apperrors.ErrUserNotFound
	}

	if !hash.VerifyPassword(user.PasswordHash, password) {
		return apperrors.ErrInvalidCredentials
	}

	if err := u.userRepository.SoftDeleteUser(ctx, userID); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	if err := u.tokenRepository.DeleteRefreshToken(ctx, userID); err != nil {
		return fmt.Errorf("failed to delete refresh token: %w", err)
	}

	userDeletedEvent := authevents.UserAccountDeletionEvent{
		Email: user.Email,
		Phone: user.Phone,
		UserID: user.ID,
	}

	if err := u.eventPublisher.PublishUserAccountDeletion(ctx, userDeletedEvent); err != nil {
		// no need to fail the process, the logging will be done in the message broker
	}

	return nil
}
