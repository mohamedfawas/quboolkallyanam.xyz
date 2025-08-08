package admin

import (
	"context"
	"fmt"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apperrors"
	authevents "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/events/auth"
)


func (u *adminUsecase) BlockOrUnblockUser(ctx context.Context, field string, value string, shouldBlock bool) error {
	user, err := u.userRepository.GetUser(ctx, field, value)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	if user == nil {
		return apperrors.ErrUserNotFound
	}

	if user.IsBlocked == shouldBlock {
		return apperrors.ErrUserAlreadyBlocked
	}

	err = u.userRepository.BlockOrUnblockUser(ctx, field, value, shouldBlock)
	if err != nil {
		return fmt.Errorf("failed to block user: %w", err)
	}

	adminBlockedUserEvent := authevents.AdminBlockedUserEvent{
		UserID: user.ID,
		Email: user.Email,
		Phone: user.Phone,
		ShouldBlock: shouldBlock,
	}

	if err := u.eventPublisher.PublishAdminBlockedUser(ctx, adminBlockedUserEvent); err != nil {
		// No need to fail the process, the logging will be done in the message broker
	}

	return nil
}