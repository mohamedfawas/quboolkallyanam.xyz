package user

import (
	"context"
	"fmt"

	constants "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	errors "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/errors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/logger"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/security/hash"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/timeutil"
)

func (u *userUseCase) UserAccountDelete(ctx context.Context, userID string, password string) error {
	user, err := u.userRepository.GetUser(ctx, "id", userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	if user == nil {
		return errors.ErrUserNotFound
	}

	if !hash.ComparePassword(user.PasswordHash, password) {
		return errors.ErrInvalidCredentials
	}

	if err := u.userRepository.SoftDeleteUser(ctx, userID); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	if err := u.tokenRepository.DeleteRefreshToken(ctx, userID); err != nil {
		return fmt.Errorf("failed to delete refresh token: %w", err)
	}

	now := timeutil.NowIST()
	if u.messageBroker != nil {
		userDeletedEvent := map[string]interface{}{
			constants.UserID:    userID,
			constants.EventType: constants.EventUserDeleted,
			constants.Timestamp: now,
		}
		if err := u.messageBroker.Publish(constants.EventUserDeleted, userDeletedEvent); err != nil {
			return fmt.Errorf("failed to publish user deleted event: %w", err)
		}
	}

	logger.Log.Info("User account deleted successfully : ", "user_id", userID)
	return nil
}
