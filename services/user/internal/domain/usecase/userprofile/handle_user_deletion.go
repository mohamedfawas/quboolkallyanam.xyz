package userprofile

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (u *userProfileUsecase) HandleUserDeletion(ctx context.Context, userID uuid.UUID) error {
	
	exists, err := u.userProfileRepository.ProfileExists(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to check if user profile exists: %w", err)
	}

	if exists {
		if err := u.userProfileRepository.DeleteUserProfile(ctx, userID); err != nil {
			return fmt.Errorf("failed to delete user profile: %w", err)
		}
		return nil
	}

	return nil // if user is already deleted no need or returning error
}