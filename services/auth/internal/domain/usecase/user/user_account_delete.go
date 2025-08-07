package user

import (
	"context"
	"fmt"

	"time"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apperrors"
	constants "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
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

	if u.messageBroker != nil {
		userDeletedEvent := map[string]interface{}{
			constants.UserID:    userID,
			constants.EventType: constants.EventUserDeleted,
			constants.Timestamp: time.Now().UTC(),
		}
		if err := u.messageBroker.Publish(constants.EventUserDeleted, userDeletedEvent); err != nil {
			// No need to fail the process, the logging will be done in the message broker
		}
	}

	return nil
}
